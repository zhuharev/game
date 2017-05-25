package lib

import (
	"log"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/middleware"
)

func handleUsers(c *middleware.Context) {
	var (
		offsetID    = c.QueryInt("offset_id")
		itemsInPage = c.QueryInt("items_per_page")
	)
	if itemsInPage == 0 {
		itemsInPage = 10
	}

	users, err := models.UserFind(offsetID, itemsInPage)
	if err != nil {
		log.Println(err)
		c.JSON(200, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var ownerIds []int64
	for _, user := range users {
		ownerIds = append(ownerIds, user.Id)
	}
	buildings, err := models.BuildingsGetByOwners(ownerIds)
	if err != nil {
		log.Println(err)
		c.JSON(200, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	for i, user := range users {
		for _, building := range buildings {
			if building.OwnerId == user.Id {
				users[i].Goods.Count++
				users[i].Goods.Buildings = append(users[i].Goods.Buildings, building)
			}
		}
	}

	c.JSON(200, map[string]interface{}{
		"users": users,
	})
}

func handleUser(c *middleware.Context) {
	user, err := models.UserGetPublicInfo(c.ParamsInt64(":id"))
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if c.QueryBool("extend") {
		blds, err := models.GetUserBuildings(user.Id)
		if err != nil {
			c.JSON(200, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		user.Goods.Buildings = blds
		user.Goods.Count = len(blds)
	}
	c.JSON(200, user)
	return
}

func handleUserSex(c *middleware.Context) {

	u, err := models.GetUserByToken(c.Query("token"))
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
	if u == nil {
		c.JSON(200, "Ошибка авторизации")
		return
	}

	var (
		sex = c.QueryInt("sex")
	)

	err = models.UsersSetSex(u.Id, sex)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, "ok")
}
