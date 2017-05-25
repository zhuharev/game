package lib

import (
	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/middleware"
)

func handleFCMToken(c *middleware.Context) {

	u, err := models.GetUserByToken(c.Query("token"))
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
	if u == nil {
		c.JSON(200, "Ошибка авторизации")
		return
	}

	err = models.UsersSetFCMToken(u.Id, c.Query("fcm_token"))
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
}
