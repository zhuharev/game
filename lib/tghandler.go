package lib

import (
	"fmt"
	"log"
	"strings"

	"github.com/Unknwon/com"
	humanize "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/fcm"
	"github.com/zhuharev/game/modules/tgbot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func tgHandler(msg *tgbotapi.Message) error {
	log.Println("Handle message", msg.Chat.ID, msg.Command(), msg.CommandArguments())
	switch msg.Command() {
	case "send":
		args := strings.Split(msg.CommandArguments(), " ")
		if len(args) < 2 {
			break
		}
		user, err := models.UserGet(com.StrTo(args[0]).MustInt64())
		if err != nil {
			log.Println(err)
			break
		}
		if user.FCMToken != "" {
			err = fcm.Send(user.FCMToken, map[string]string{"text": strings.Join(args[1:], " "), "type": "test"})
			if err != nil {
				log.Println(err)
			}
		} else {
			err = tgbot.Send(msg.Chat.ID, "У пользователя не зарегистрировано устройство (возможно у него старая версия приложения)")
			if err != nil {
				log.Println(err)
			}
		}
	case "server":
		us, _ := disk.Usage("/")
		v, _ := mem.VirtualMemory()
		cnt, _ := models.BuildingWithOwnerCount()
		ucnt, _ := models.UsersCount()
		stat := fmt.Sprintf(`Память: использовано %v из %v (%.2f%%)
Диск: использовано %s из %s (%.2f%%)
Зданий захвачено: %d
Пользователей: %d`, humanize.Bytes(v.Used), humanize.Bytes(v.Total), v.UsedPercent,
			humanize.Bytes(us.Used), humanize.Bytes(us.Total), us.UsedPercent, cnt, ucnt)
		err := tgbot.Send(msg.Chat.ID, stat)
		if err != nil {
			log.Println(err)
		}
	case "price":
		if arr := strings.Split(msg.CommandArguments(), " "); len(arr) == 2 {
			err := updatePrice(arr[0], com.StrTo(arr[1]).MustInt64())
			if err != nil {
				err = tgbot.Send(msg.Chat.ID, fmt.Sprintf("Ошибка: %s", err))
				if err != nil {
					log.Println(err)
				}
				break
			}
			pr, err := models.GetPrices()
			if err != nil {
				log.Println(err)
				break
			}
			err = tgbot.Send(msg.Chat.ID, fmt.Sprintf("Цены обновлены, новые цены:\n%s", pr))
			if err != nil {
				log.Println(err)
			}
		}
	case "prices":
		pr, err := models.GetPrices()
		if err != nil {
			log.Println(err)
			break
		}
		err = tgbot.Send(msg.Chat.ID, fmt.Sprintf("Цены:\n%s", pr))
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func updatePrice(name string, value int64) error {
	return models.SetPrice(name, value)
}
