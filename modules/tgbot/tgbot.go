package tgbot

import (
	"log"

	"github.com/zhuharev/game/modules/setting"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	bot *tgbotapi.BotAPI
)

func NewContext(handler func(*tgbotapi.Message) error) (err error) {
	bot, err = tgbotapi.NewBotAPI(setting.App.Telegram.Key)
	if err != nil {
		return
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			err = handler(update.Message)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	return nil
}

func Send(target int64, message string) error {
	msg := tgbotapi.NewMessage(target, message)
	//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
	//	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Text")))
	_, err := bot.Send(msg)
	return err
}

// func ShowSettingsKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
// 	markup := tg.NewInlineKeyboardMarkup(
// 		tg.NewInlineKeyboardRow(
// 			tg.NewInlineKeyboardButtonData(T("button_language", map[string]interface{}{"Flag": T("language_flag")}), "language menu"),
// 		),
// 		tg.NewInlineKeyboardRow(
// 			tg.NewInlineKeyboardButtonData(T("button_resources"), "resources menu"),
// 		),
// 		tg.NewInlineKeyboardRow(
// 			tg.NewInlineKeyboardButtonData(T("button_ratings"), "ratings menu"),
// 		),
// 		tg.NewInlineKeyboardRow(
// 			tg.NewInlineKeyboardButtonData(T("button_blacklist"), "blacklist menu"),
// 		),
// 		tg.NewInlineKeyboardRow(
// 			tg.NewInlineKeyboardButtonData(T("button_whitelist"), "whitelist menu"),
// 		),
// 	)
//
// 	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
// 	if _, err := bot.Send(edit); err != nil {
// 		log.Println(err.Error())
// 	}
// }
