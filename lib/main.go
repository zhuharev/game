package lib

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/binding"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/tile38"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
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

	app := iris.New()
	// output startup banner and error logs on os.Stdout
	app.Adapt(iris.DevLogger())
	// set the router, you can choose gorillamux too
	app.Adapt(httprouter.New())

	api := app.Party("/api/v1")
	api.Get("/buildings", handleBuildings)
	api.Get("/users/me", me)
	api.Get("/auth", handleAuth)
	api.Get("/games/new", handleNewGame)
	api.Get("/games/check", handleCheck)
	api.Get("/location/update", func(ctx *iris.Context) {
		user, err := models.GetUserFromCtx(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}

		arr := strings.Split(ctx.FormValue("location"), ",")
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

	app.Get("/", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, iris.Map{"name": "iris"})
	})

	app.Listen(":7000")
}
