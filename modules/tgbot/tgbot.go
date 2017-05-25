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
	_, err := bot.Send(msg)
	return err
}
