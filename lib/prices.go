package lib

import (
	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/middleware"
)

func handlePrices(c *middleware.Context) {
	prices, err := models.GetPrices()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, prices)
}
