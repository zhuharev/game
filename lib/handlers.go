package lib

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/mholt/binding"
	"gopkg.in/kataras/iris.v6"
	"strconv"
	"strings"

	"github.com/zhuharev/game/models"
)

func handleNewGame(ctx *iris.Context) {
	user, err := models.GetUserFromCtx(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	buildingId := com.StrTo(ctx.FormValue("building_id")).MustInt64()
	game, err := models.NewGame(user.Id, buildingId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(200, game)
}

func handleCheck(ctx *iris.Context) {
	user, err := models.GetUserFromCtx(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	answer := com.StrTo(ctx.FormValue("answer")).MustInt()
	bulls, cows, err := models.Check(user, com.StrTo(ctx.FormValue("game_id")).MustInt64(), answer)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(200, iris.Map{
		"answer": answer,
		"bulls":  bulls,
		"cows":   cows,
	})
}

func handleError(ctx *iris.Context, err error) {
	if err != nil {
		ctx.JSON(200, iris.Map{
			"error": err.Error(),
		})
	}
}

func handleBuildings(ctx *iris.Context) {

	cntr := new(Center)

	errs := binding.Bind(ctx.Request, cntr)
	if errs.Has("") {
		fmt.Println(errs)
	}

	arr := strings.Split(cntr.LongLat, ",")
	if len(arr) != 2 {
		fmt.Println("Not 2")
	}

	lat, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		fmt.Println(err)
	}

	lon, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		fmt.Println(err)
	}

	buildings, err := models.Nearby(lat, lon)
	if err != nil {
		fmt.Println(err)
	}
	var ids []int64

	for _, b := range buildings {
		ids = append(ids, b.Id)
	}
	buildings, err = models.FindBuildings(ids)
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(iris.StatusOK, buildings)
}

func handleAuth(c *iris.Context) {
	var aform models.AuthForm
	err := c.ReadForm(&aform)
	if err != nil {
		color.Red("%s", err)
		c.JSON(200, err.Error())
		return
	}

	u, err := models.AuthUser(aform)
	if err != nil {
		color.Red("%s", err)
		c.JSON(200, err.Error())
		return
	}

	c.JSON(200, struct {
		Id    int64  `json:"user_id"`
		Token string `json:"token"`
	}{
		u.Id,
		u.Token,
	})
}

func me(c *iris.Context) {
	token := c.FormValue("token")
	if token == "" {
		c.JSON(200, "token is nil")
	}
	u, err := models.GetUserByToken(token)
	if err != nil {
		c.JSON(200, err.Error())
		return
	}

	buildings, err := models.GetUserBuildings(u.Id)
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
	u.Goods.Buildings = buildings
	u.Goods.Count = len(buildings)

	u.Token = token
	c.JSON(200, u)
}
