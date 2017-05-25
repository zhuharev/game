package setting

import (
	"gopkg.in/ini.v1"
)

var (
	iniFile *ini.File
)

func NewContext() (err error) {
	err = ini.MapToWithMapper(&App, ini.TitleUnderscore, "conf/app.ini")

	return
}

var (
	App struct {
		Fcm struct {
			Key string
		}

		Telegram struct {
			Key string
		}
	}
)
