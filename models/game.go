package models

import (
	"fmt"
	"log"
	"time"

	bac "github.com/zhuharev/game/modules/bulls"
	"github.com/zhuharev/game/modules/fcm"
	"github.com/zhuharev/game/modules/fixdb"
	"github.com/zhuharev/game/modules/tgbot"
)

const (
	MaxDistance = 1.0 // in km
)

type Game struct {
	Id      int64 `json:"id"`
	BrickID int64 `json:"brick_id" xorm:"brick_id"`

	Secret int `json:"-"`

	UserId     int64 `json:"user_id"`
	BuildingId int64 `json:"building_id"`
	Steps      int   `json:"-"`

	Status int `json:"status"`

	Step           int `json:"-" xorm:"step"`
	PinChangedStep int `json:"-"`

	Updated time.Time `xorm:"updated" json:"-"`
	Created time.Time `xorm:"created" json:"-"`
}

// NewGame create game.
// If building not exists in sql database, it try find in fixdb.
func NewGame(userID, buildingID, brickID int64) (*Game, error) {
	user, err := UserGet(userID)
	if err != nil {
		return nil, err
	}
	building, err := BuildingGet(buildingID)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		var coords []float64
		var area int64
		coords, area, err = fixdb.Get(buildingID)
		if err != nil {
			return nil, err
		}
		if coords != nil {
			building = &Building{
				Id:     buildingID,
				Lat:    coords[0],
				Long:   coords[1],
				Area:   int(area),
				Profit: 1,
			}
			err = createBuilding(building)
			if err != nil {
				return nil, err
			}
		}
	}

	if building.OwnerId == user.Id {
		return nil, fmt.Errorf("Already owner")
	}

	// check if building blocked
	if !building.Blocked.IsZero() && time.Since(building.Blocked) < building.BlockedDuration {
		return nil, fmt.Errorf("Здание заблокировано, осталось %.0fм", (building.BlockedDuration - time.Since(building.Blocked)).Minutes())
	}

	if building.OwnerId != 0 && building.NotificationEnabled {
		var oldOwner *User
		oldOwner, err = userGet(building.OwnerId, true)
		if err != nil {
			return nil, err
		}

		if oldOwner.FCMToken != "" {
			err = fcm.Send(oldOwner.FCMToken, map[string]string{"text": "Ваше здание захватывают!",
				"type": "seizure", "building_id": fmt.Sprint(building.Id), "brick_id": fmt.Sprint(brickID)})
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("Ignore notification, fcm token not set", oldOwner.Id)
		}
	} else {
		log.Println("Ignore notification, building enabled = ", building.NotificationEnabled)
	}

	/*	if dist := distance(
			user.Lat,
			user.Lon,
			building.Lat,
			building.Long,
		); dist > MaxDistance {
			return nil, fmt.Errorf("Подойдите ближе!")
		}*/

	if brickID < 1 {
		brickID = 1
	}

	game := &Game{
		UserId:         user.Id,
		BuildingId:     buildingID,
		Secret:         toInt(genNumber()),
		BrickID:        brickID,
		Step:           1,
		PinChangedStep: 0,
	}

	_, err = db.Insert(game)
	if err != nil {
		return nil, err
	}

	go func() {
		err = tgbot.Send(166935911, fmt.Sprintf("%d", game.Secret))
		if err != nil {
			log.Println(err)
		}
		err = tgbot.Send(102710272, fmt.Sprintf("%d", game.Secret))
		if err != nil {
			log.Println(err)
		}

	}()

	return game, nil
}

func getGame(id int64) (*Game, error) {
	game := new(Game)
	has, err := db.Id(id).Get(game)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return game, nil
}

func Check(u *User, gameId int64, answer int, step int) (bulls int, cows int, highlight []int, nextGameID int64, message string, armor int64, err error) {
	game, err := getGame(gameId)
	if err != nil {
		return
	}
	if game.UserId != u.Id {
		err = fmt.Errorf("Forbiten")
		return
	}

	bulls, cows, err = bac.BullsAndCows(toByte(game.Secret), toByte(answer))
	if err != nil {
		return
	}

	build, err := BuildingGet(game.BuildingId)
	if err != nil {
		return
	}

	armor = build.Armor

	if game.PinChangedStep > 0 && game.PinChangedStep == step {
		log.Println("[PIN CHANGE]", game.PinChangedStep, step)
		message = "Владелец здания изменил пинкод!"
	} else {
		log.Println("[PIN NOT CHANGE]", game.PinChangedStep, step)
	}

	if !build.Blocked.IsZero() && time.Since(build.Blocked) < build.BlockedDuration {
		err = fmt.Errorf("Здание заблокировано, осталось %.0fм", (build.BlockedDuration - time.Since(build.Blocked)).Minutes())
		return
	}

	highlight = bac.Highlight(toByte(game.Secret), toByte(answer), build.Armor)

	fmt.Printf("Check game(%d): answer=%d (b:%d, c:%d), hilights: %v\n", game.Secret, answer, bulls, cows, highlight)

	_, err = db.Exec("update building set step = ? where id = ?", step, game.BuildingId)
	if err != nil {
		log.Println(err)
		return
	}

	// win
	if bulls == 4 {
		game.Status = 1
		_, err = db.Id(game.Id).Update(game)
		if err != nil {
			return
		}

		if NeedNextGame(game.BrickID, build.Armor) {
			var nextGame *Game
			nextGame, err = NewGame(u.Id, build.Id, game.BrickID+1)
			if err != nil {
				return
			}
			nextGameID = nextGame.Id
		}

		_, err = buildGetFromSQL(build.Id)
		if err != nil && err != ErrNotFound {
			log.Println(err)
			return
		} else if err == ErrNotFound {
			// создаем здание в sql
			err = createBuilding(build)
			if err != nil {
				log.Println(err)
				return
			}
		}

		// totaly win
		if nextGameID == 0 {

			_, err = db.Exec("update building set owner_id = ?, armor = 1, notification_enabled = ? where id = ?", u.Id, true, game.BuildingId)
			if err != nil {
				log.Println(err)
				return
			}

			// пересчитываем прибыль всех
			// TODO пересчитавать только бывшего владельца и нового
			_, err = db.Exec("update user set profit = (select sum(profit) from building where owner_id = user.id)")
			if err != nil {
				log.Println(err)
				return
			}
		}

	}

	return
}

func NeedNextGame(currentBrick int64, buildingArmorLevel int64) bool {
	if buildingArmorLevel <= 4 {
		return false
	}
	return buildingArmorLevel-3 > currentBrick
}
