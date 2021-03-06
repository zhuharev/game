package lib

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/bloblog"
	"github.com/zhuharev/game/modules/fixdb"
	"github.com/zhuharev/game/modules/middleware"
	"github.com/zhuharev/game/modules/nearbydb"
	"github.com/zhuharev/game/modules/setting"
	"github.com/zhuharev/game/modules/tgbot"
)

// Run starts web server
func Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			os.Exit(0)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	err := setting.NewContext()
	if err != nil {
		log.Fatalln(err)
	}

	err = fixdb.NewContext()
	if err != nil {
		log.Fatalln(err)
	}
	err = bloblog.NewContext()
	if err != nil {
		log.Fatalln(err)
	}

	models.SetDb()
	err = nearbydb.NewContext()
	if err != nil {
		log.Fatalln(err)
	}

	err = tgbot.NewContext(tgHandler)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(setting.App.Fcm.Key)

	time.Sleep(1 * time.Second)

	go func() {
		tick := time.NewTicker(15 * time.Minute)
		for range tick.C {
			err := models.IncreaseBalance()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	m := macaron.Classic()
	m.Use(macaron.Static("static"))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		IndentJSON: true,
	}))
	m.Use(middleware.Contexter())

	m.Group("/api/v1", func() {
		m.Get("/prices", handlePrices)
		m.Get("/buildings", handleBuildings)
		m.Get("/buildings/:id", Building)
		m.Get("/buildings/:id/upgrade", handleUpgrade)
		m.Get("/buildings/:id/downgrade", handleDowngrade)
		m.Get("/buildings/:id/block", handleBlock)
		m.Get("/buildings/:id/notify", handleNotify)
		m.Get("/buildings/:id/pin", handlePin)
		m.Get("/user", me)
		m.Group("/users", func() {
			m.Get("/:id", handleUser)
			m.Get("/:id/sex", handleUserSex)
			m.Get("/:id/vk_avatar", handleVkAvar)
		})
		m.Get("/users", handleUsers)

		m.Get("/auth", handleAuth)
		m.Get("/games/new", handleNewGame)
		m.Get("/games/check", handleCheck)
		m.Get("/location/update", func(ctx *middleware.Context) {
			user, err := models.GetUserFromCtx(ctx)
			if err != nil {
				handleError(ctx, err)
				return
			}

			arr := strings.Split(ctx.Query("location"), ",")
			if len(arr) != 2 {
				if err != nil {
					handleError(ctx, fmt.Errorf("not 2"))
					return
				}
			}

			lat, err := strconv.ParseFloat(arr[0], 64)
			if err != nil {
				handleError(ctx, err)
				return
			}

			lon, err := strconv.ParseFloat(arr[1], 64)
			if err != nil {
				handleError(ctx, err)
				return
			}

			err = models.SetLocation(user.Id, lat, lon)
			if err != nil {
				handleError(ctx, err)
				return
			}
			ctx.JSON(200, "ok")
		})
		m.Post("/user/avatar", handleUpload)
		m.Get("/user/fcm", handleFCMToken)
	})

	m.Get("/", func(ctx *middleware.Context) {
		ctx.JSON(200, map[string]interface{}{"name": "Juctvalk"})
	})

	m.Get("/images/:id", handleAvatar)

	m.Run(7000)
}
