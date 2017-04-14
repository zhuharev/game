package main

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/mholt/binding"
	"gopkg.in/kataras/iris.v6"
	"strconv"
	"strings"
)

func handleNewGame(ctx *iris.Context) {
	user, err := getUserFromCtx(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	buildingId := com.StrTo(ctx.FormValue("building_id")).MustInt64()
	game, err := newGame(user.Id, buildingId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(200, game)
}

func handleCheck(ctx *iris.Context) {
	user, err := getUserFromCtx(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	answer := com.StrTo(ctx.FormValue("answer")).MustInt()
	bulls, cows, err := check(user, com.StrTo(ctx.FormValue("game_id")).MustInt64(), answer)
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

	buildings, err := nearby(lat, lon)
	if err != nil {
		fmt.Println(err)
	}
	var ids []int64

	for _, b := range buildings {
		ids = append(ids, b.Id)
	}
	buildings, err = findBuildings(ids)
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(iris.StatusOK, buildings)
}
