// Settings commands for bot admins
// /autodelete [minutes] - set auto delete time (0 to disable)
// /multifilter [on/off] - toggle multi filter

package plugins

import (
	"fmt"
	"strconv"

	"github.com/Jisin0/Go-Filter-Bot/utils/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// isAdmin checks if user is in ADMINS list
func isAdmin(userId int64) bool {
	for _, id := range config.Admins {
		if id == userId {
			return true
		}
	}
	return false
}

// /autodelete [minutes] — 0 means disable
// Example: /autodelete 5  -> messages delete in 5 minutes
// Example: /autodelete 0  -> auto delete off
func SetAutoDelete(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveMessage.From

	// Sirf admin use kar sakta hai
	if !isAdmin(user.Id) {
		_, err := ctx.EffectiveMessage.Reply(bot, "❌ Sirf bot admins yeh command use kar sakte hain!", nil)
		return err
	}

	args := ctx.Args()

	// Agar koi argument nahi diya
	if len(args) < 2 {
		_, err := ctx.EffectiveMessage.Reply(bot,
			fmt.Sprintf("⚙️ <b>Auto Delete Setting</b>\n\nAbhi: <b>%v minutes</b>\n\nUsage: <code>/autodelete [minutes]</code>\nExample: <code>/autodelete 5</code>\nDisable ke liye: <code>/autodelete 0</code>",
				AutoDelete/60),
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err
	}

	// Minutes parse karo
	minutes, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || minutes < 0 {
		_, err = ctx.EffectiveMessage.Reply(bot, "❌ Galat value! Sirf number dalo. Example: <code>/autodelete 5</code>",
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err
	}

	// Value update karo (seconds mein store hoti hai)
	AutoDelete = minutes * 60

	var replyText string
	if minutes == 0 {
		replyText = "✅ Auto Delete <b>band</b> kar diya gaya!"
	} else {
		replyText = fmt.Sprintf("✅ Auto Delete set ho gaya: <b>%v minutes</b>", minutes)
	}

	_, err = ctx.EffectiveMessage.Reply(bot, replyText, &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	return err
}

// /multifilter [on/off] — toggle multi filter
// Example: /multifilter on
// Example: /multifilter off
func SetMultiFilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveMessage.From

	// Sirf admin use kar sakta hai
	if !isAdmin(user.Id) {
		_, err := ctx.EffectiveMessage.Reply(bot, "❌ Sirf bot admins yeh command use kar sakte hain!", nil)
		return err
	}

	args := ctx.Args()

	// Agar koi argument nahi diya — current status batao
	if len(args) < 2 {
		status := "❌ Off"
		if config.MultiFilter {
			status = "✅ On"
		}
		_, err := ctx.EffectiveMessage.Reply(bot,
			fmt.Sprintf("⚙️ <b>Multi Filter Setting</b>\n\nAbhi: <b>%v</b>\n\nUsage: <code>/multifilter on</code> ya <code>/multifilter off</code>", status),
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err
	}

	switch args[1] {
	case "on", "true", "enable", "1":
		config.MultiFilter = true
		_, err := ctx.EffectiveMessage.Reply(bot, "✅ Multi Filter <b>on</b> kar diya gaya!", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err

	case "off", "false", "disable", "0":
		config.MultiFilter = false
		_, err := ctx.EffectiveMessage.Reply(bot, "✅ Multi Filter <b>off</b> kar diya gaya!", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err

	default:
		_, err := ctx.EffectiveMessage.Reply(bot, "❌ Galat value! <code>/multifilter on</code> ya <code>/multifilter off</code> likho",
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err
	}
}
