package lib

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/mholt/binding"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/middleware"
	"github.com/zhuharev/game/modules/vk"
)

func handleNewGame(ctx *middleware.Context) {
	user, err := models.GetUserFromCtx(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	buildingID := ctx.QueryInt64("building_id")
	game, err := models.NewGame(user.Id, buildingID, ctx.QueryInt64("brick_id"))
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(200, game)
}

func handleCheck(ctx *middleware.Context) {
	user, err := models.GetUserFromCtx(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	answer := com.StrTo(ctx.Query("answer")).MustInt()
	bulls, cows, highlights, nextGameID, message, armor, err := models.Check(user, ctx.QueryInt64("game_id"), answer, ctx.QueryInt("step"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	res := map[string]interface{}{
		"answer":     answer,
		"bulls":      bulls,
		"cows":       cows,
		"highlights": highlights,
		"next_game":  nextGameID,
		"armor":      armor,
	}
	if message != "" {
		res["message"] = message
	}
	log.Println("[RES]", res)
	ctx.JSON(200, res)
}

func handleError(ctx *middleware.Context, err error) {
	if err != nil {
		ctx.JSON(200, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

type center struct {
	LongLat string
	lon     float64
	lat     float64
	parsed  bool
}

func (c *center) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.LongLat: "center",
	}
}

func (c *center) parse() {
	if c.parsed {
		return
	}
	defer func() { c.parsed = true }()

	arr := strings.Split(c.LongLat, ",")
	if len(arr) != 2 {
		fmt.Println("Not 2")
		return
	}

	c.lat, _ = strconv.ParseFloat(arr[0], 64)
	c.lon, _ = strconv.ParseFloat(arr[1], 64)
}

func (c *center) Lat() float64 {
	c.parse()
	return c.lat
}

func (c *center) Lon() float64 {
	c.parse()
	return c.lon
}

func handleBuildings(ctx *middleware.Context) {

	cntr := new(center)

	errs := binding.Bind(ctx.Context.Req.Request, cntr)
	if errs.Has("") {
		fmt.Println(errs)
	}

	buildings, err := models.Nearby(cntr.Lat(), cntr.Lon())
	if err != nil {
		fmt.Println(err)
	}
	var ids []int64

	for _, b := range buildings {
		ids = append(ids, b.Id)
	}
	buildingsFromDb, err := models.FindBuildings(ids)
	if err != nil {
		fmt.Println(err)
	}
	// show buildins if it not exsist in sql
	for i, b := range buildings {
		for _, bFromDb := range buildingsFromDb {
			if bFromDb.Id == b.Id {
				buildings[i] = bFromDb
			}
		}
	}

	var ownerIds []int64
	for _, b := range buildings {
		if b.OwnerId != 0 {
			ownerIds = append(ownerIds, b.OwnerId)
		}
	}

	users, err := models.UserFindByIds(ownerIds)
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(200, map[string]interface{}{
		"buildings": buildings,
		"users":     users,
	})
}

func handleAuth(c *middleware.Context) {
	var aform models.AuthForm
	err := c.ReadForm(&aform)
	if err != nil {
		color.Red("%s", err)
		c.JSON(200, err.Error())
		return
	}

	user, err := vk.CheckToken(aform.VkToken)
	if err != nil {
		if err != nil {
			color.Red("%s", err)
			c.JSON(200, err.Error())
			return
		}
	}
	aform.VkId = int64(user.Id)
	aform.FirstName = user.FirstName
	aform.LastName = user.LastName
	aform.AvatarURL = user.Photo200

	u, err := models.AuthUser(aform)
	if err != nil {
		color.Red("%s", err)
		c.JSON(200, err.Error())
		return
	}

	c.JSON(200, struct {
		ID    int64  `json:"id"`
		Token string `json:"token"`
	}{
		u.Id,
		u.Token,
	})
}

func me(c *middleware.Context) {
	token := c.Query("token")
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
