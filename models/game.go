package models

import (
	"fmt"
	"log"
	"time"

	bac "github.com/zhuharev/game/modules/bulls"
	"github.com/zhuharev/game/modules/fixdb"
)

const (
	MaxDistance = 1.0 // in km
)

type Game struct {
	Id      int64 `json:"id"`
	BrickID int64 `json:"brick_id"`

	Secret int `json:"-"`

	UserId     int64 `json:"user_id"`
	BuildingId int64 `json:"building_id"`
	Steps      int   `json:"-"`

	Status int `json:"status"`

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
		coords, err = fixdb.Get(buildingID)
		if err != nil {
			return nil, err
		}
		if coords != nil {
			building = &Building{
				Id:     buildingID,
				Lat:    coords[0],
				Long:   coords[1],
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

	/*	if dist := distance(
			user.Lat,
			user.Lon,
			building.Lat,
			building.Long,
		); dist > MaxDistance {
			return nil, fmt.Errorf("Подойдите ближе!")
		}*/

	game := &Game{
		UserId:     user.Id,
		BuildingId: buildingID,
		Secret:     toInt(genNumber()),
		BrickID:    brickID,
	}

	_, err = db.Insert(game)
	if err != nil {
		return nil, err
	}

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

func Check(u *User, gameId int64, answer int) (bulls int, cows int, highlight []int, err error) {
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

	highlight = bac.Highlight(toByte(game.Secret), toByte(answer), build.Armor)

	fmt.Printf("Check game(%d): answer=%d (b:%d, c:%d), hilights: %v\n", game.Secret, answer, bulls, cows, highlight)

	// win
	if bulls == 4 {
		game.Status = 1
		_, err = db.Id(game.Id).Update(game)
		if err != nil {
			return
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
		if build.OwnerId != 0 {
			// уменьшаем прибыль предыдущего владельца
			_, err = db.Exec("update user set profit = profit - ? where id = (select owner_id from building where id = ?)", build.Profit, build.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}
		_, err = db.Exec("update building set owner_id = ? where id = ?", u.Id, game.BuildingId)
		if err != nil {
			log.Println(err)
			return
		}
		// увеличиваем прибыль нового владельца
		_, err = db.Exec("update user set profit = profit + ? where id = ?", build.Profit, u.Id)
		if err != nil {
			log.Println(err)
			return
		}
	}

	return
}
