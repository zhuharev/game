package lib

import (
	"github.com/zhuharev/game/models"
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
		user, err := models.UserGet(building.OwnerId)
		if err != nil {
			handleError(c, err)
			return
		}
		building.OwnerName = user.FullName
	}

	c.JSON(200, building)
}
