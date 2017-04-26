package lib

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/binding"
	"gopkg.in/macaron.v1"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/middleware"
	"github.com/zhuharev/game/modules/tile38"
)

type Center struct {
	LongLat string
}

func (c *Center) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.LongLat: "center",
	}
}

func Run() {
	go func() {
		err := tile38.StartTileServer()
		if err != nil {
			panic(err)
		}
	}()

	models.SetDb()

	time.Sleep(1 * time.Second)

	go func() {
		tick := time.NewTicker(1 * time.Minute)
		for range tick.C {
			err := models.IncreaseBalance()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Use(middleware.Contexter())

	m.Group("/api/v1", func() {
		m.Get("/buildings", handleBuildings)
		m.Get("/user", me)
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
	})

	m.Get("/", func(ctx *middleware.Context) {
		ctx.JSON(200, map[string]interface{}{"name": "iris"})
	})

	m.Run(7000)
}
