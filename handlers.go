package main

import (
	"github.com/Unknwon/com"
	"gopkg.in/kataras/iris.v6"
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
