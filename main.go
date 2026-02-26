package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Jisin0/Go-Filter-Bot/plugins"
	"github.com/Jisin0/Go-Filter-Bot/utils/autodelete"
	"github.com/Jisin0/Go-Filter-Bot/utils/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println(string(debug.Stack()))
		}
	}()

	// ✅ STEP 1: Pehle Web Server Start Karo
	// Render/Koyeb ko healthy build ke liye HTTP server chahiye
	serverReady := make(chan bool, 1)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "Waku Waku - Bot is Running!")
		})

		// Server ready signal bhejo
		serverReady <- true

		port := config.Port
		if port == "" {
			port = "8080"
		}

		fmt.Println("🌐 Web server starting on port: " + port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Println("Web server error: " + err.Error())
		}
	}()

	// ✅ STEP 2: Web Server ke ready hone ka wait karo (max 5 seconds)
	select {
	case <-serverReady:
		fmt.Println("✅ Web server is ready!")
		// Thoda extra time do server ko properly bind karne ke liye
		time.Sleep(1 * time.Second)
	case <-time.After(5 * time.Second):
		fmt.Println("⚠️ Web server timeout, proceeding anyway...")
	}

	// ✅ STEP 3: Ab Bot Token Check Karo
	if config.BotToken == "" {
		panic("❌ Exiting Because No BOT_TOKEN Provided :(")
	}

	fmt.Println("🤖 Starting Telegram Bot...")

	// ✅ STEP 4: Bot Create Karo
	b, err := gotgbot.NewBot(config.BotToken, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	})
	if err != nil {
		panic("❌ Failed to create new bot: " + err.Error())
	}

	// ✅ STEP 5: Check karo koi aur instance toh nahi chal raha
	_, err = b.GetUpdates(&gotgbot.GetUpdatesOpts{})
	if err != nil {
		fmt.Println("⏳ Waiting 10s because: " + err.Error())
		time.Sleep(10 * time.Second)
	}

	// ✅ STEP 6: Updater Setup Karo
	updater := ext.NewUpdater(plugins.Dispatcher, &ext.UpdaterOpts{})

	// ✅ STEP 7: Polling Start Karo
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			AllowedUpdates: []string{
				"message",
				"callback_query",
				"channel_post",
				"inline_query",
				"chosen_inline_result",
				"chat_member",
				"my_chat_member",
			},
		},
	})
	if err != nil {
		panic("❌ Failed to start polling: " + err.Error())
	}

	fmt.Printf("✅ @%s Started Successfully!\n", b.User.Username)

	// ✅ STEP 8: AutoDelete Feature
	if plugins.AutoDelete > 0 {
		go autodelete.RunAutodel(b)
	}

	// ✅ STEP 9: Bot ko idle rakhO — updates aate rahen
	updater.Idle()
}
