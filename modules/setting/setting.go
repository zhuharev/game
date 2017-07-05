package setting

import (
	"log"

	"gopkg.in/ini.v1"

	"github.com/zhuharev/vk"
)

var (
	iniFile *ini.File
)

// NewContext reads app.ini file
func NewContext() (err error) {
	err = ini.MapToWithMapper(&App, ini.TitleUnderscore, "conf/app.ini")
	//TODO check vk.AdminToken on start
	if App.Vk.AdminToken == "" {
		uri, err := vk.GetAuthURL(vk.DefaultRedirectURI, "token", App.Vk.ClientId, "offline")
		if err != nil {
			log.Println(err)
		}
		log.Println(uri)
	}
	return
}

var (
	App struct {
		App struct {
			Host string
		}

		Fcm struct {
			Key string
		}

		Telegram struct {
			Key string
		}

		Vk struct {
			AdminToken string
			ClientId   string
		}
	}
)
