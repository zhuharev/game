package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/fcm"
	"github.com/zhuharev/game/modules/middleware"
)

// Building return info about building
func Building(c *middleware.Context) {
	var (
		id = c.ParamsInt64(":id")
	)

	building, err := models.BuildingGet(id)
	if err != nil {
		handleError(c, err)
		return
	}
	if building.OwnerId != 0 {
		user, err := models.UserGetPublicInfo(building.OwnerId)
		if err != nil {
			handleError(c, err)
			return
		}
		building.OwnerName = user.FullName

		cnt, err := models.GetUserBuildingCount(user.Id)
		if err != nil {
			handleError(c, err)
			return
		}
		user.Goods.Count = int(cnt)

		c.JSON(200, models.BuildingWithOwner{Building: *building, Owner: *user})
		return
	}

	c.JSON(200, building)
}

func handleUpgrade(c *middleware.Context) {
	var (
		// armor
		//		typ        = c.Query("type")
		buildingID = c.ParamsInt64(":id")
		delta      = c.QueryInt64("levels")
	)

	if delta == 0 {
		delta = 1
	}

	log.Printf("Upgrade building %d (delta = %d)\n", buildingID, delta)

	user, err := models.GetUserFromCtx(c)
	if err != nil {
		handleError(c, err)
		return
	}

	build, err := models.BuildingGet(buildingID)
	if err != nil {
		handleError(c, err)
		return
	}

	if build.OwnerId != user.Id {
		handleError(c, fmt.Errorf("Вы не являетесь владельцем"))
		return
	}

	prices, err := models.GetPrices()
	if err != nil {
		handleError(c, err)
		return
	}

	if user.Balance.Real() < (prices.Armor * delta) {
		handleError(c, fmt.Errorf("Недостаточно денег"))
		return
	}

	err = models.DecreaseBalance(user.Id, prices.Armor*delta)
	if err != nil {
		handleError(c, err)
		return
	}

	err = models.BuildingIncArmor(buildingID, delta)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, "ok")
}

func handleDowngrade(c *middleware.Context) {
	var (
		// armor
		//		typ        = c.Query("type")
		buildingID = c.ParamsInt64(":id")
		delta      = c.QueryInt64("levels")
	)

	if delta == 0 {
		delta = 1
	}

	log.Printf("Downgrade building %d (delte = %d)\n", buildingID, delta)

	user, err := models.GetUserFromCtx(c)
	if err != nil {
		handleError(c, err)
		return
	}

	build, err := models.BuildingGet(buildingID)
	if err != nil {
		handleError(c, err)
		return
	}

	if build.Armor < 1 {
		handleError(c, fmt.Errorf("Защита уже на нуле"))
		return
	}

	if build.OwnerId == user.Id {
		handleError(c, fmt.Errorf("Вы не можете понижать защиту своего здания"))
		return
	}

	prices, err := models.GetPrices()
	if err != nil {
		handleError(c, err)
		return
	}

	if user.Balance.Real() < (prices.Armor * delta) {
		handleError(c, fmt.Errorf("Недостаточно денег"))
		return
	}

	err = models.DecreaseBalance(user.Id, prices.Armor*delta)
	if err != nil {
		handleError(c, err)
		return
	}

	err = models.BuildingDecrArmor(buildingID, delta)
	if err != nil {
		handleError(c, err)
		return
	}

	if build.OwnerId != 0 {
		var oldOwner *models.User
		oldOwner, err = models.UserGet(build.OwnerId)
		if err != nil {
			handleError(c, err)
			return
		}

		if oldOwner.FCMToken != "" {
			err = fcm.Send(oldOwner.FCMToken, map[string]string{"text": "Защита здания понижена!",
				"type": "downgrade", "building_id": fmt.Sprint(build.Id)})
			if err != nil {
				log.Println(err)
			}
		}
	}

	c.JSON(200, "ok")
}

func handleBlock(c *middleware.Context) {
	var (
		// armor
		//		typ        = c.Query("type")
		buildingID = c.ParamsInt64(":id")
	)

	log.Printf("Block building %d\n", buildingID)

	user, err := models.GetUserFromCtx(c)
	if err != nil {
		handleError(c, err)
		return
	}

	build, err := models.BuildingGet(buildingID)
	if err != nil {
		handleError(c, err)
		return
	}

	if !build.Blocked.IsZero() && build.Blocked.Add(build.BlockedDuration).After(time.Now()) {
		handleError(c, fmt.Errorf("Здание уже заблокировано"))
		return
	}

	if build.OwnerId != user.Id {
		handleError(c, fmt.Errorf("Вы не можете блокировать чужое здание"))
		return
	}

	prices, err := models.GetPrices()
	if err != nil {
		handleError(c, err)
		return
	}

	if user.Balance.Real() < prices.Block {
		handleError(c, fmt.Errorf("Недостаточно денег"))
		return
	}

	err = models.DecreaseBalance(user.Id, prices.Block)
	if err != nil {
		handleError(c, err)
		return
	}

	err = models.BuildingBlock(buildingID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{"response": 1})
}

func handleNotify(c *middleware.Context) {
	var (
		// armor
		//		typ        = c.Query("type")
		buildingID = c.ParamsInt64(":id")
	)

	log.Printf("Set notify %d\n", buildingID)

	err := models.BuildingSetNotify(buildingID, c.QueryBool("enabled"))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{"response": 1})
}

func handlePin(c *middleware.Context) {
	var (
		// armor
		//		typ        = c.Query("type")
		buildingID = c.ParamsInt64(":id")
		brickID    = c.QueryInt64("brick")
	)

	log.Printf("Set new pin %d\n", buildingID)

	err := models.BuildingRepin(buildingID, brickID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{"response": 1})
}
