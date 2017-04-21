package main

import (
	"fmt"
	"time"
)

const (
	MaxDistance = 1.0 // in km
)

type Game struct {
	Id int64 `json:"id"`

	Secret int `json:"-"`

	UserId     int64 `json:"user_id"`
	BuildingId int64 `json:"building_id"`
	Steps      int   `json:"-"`

	Status int `json:"status"`

	Updated time.Time `xorm:"updated" json:"-"`
	Created time.Time `xorm:"created" json:"-"`
}

func newGame(userId, buildingId int64) (*Game, error) {
	user, err := getUser(userId)
	if err != nil {
		return nil, err
	}
	building, err := getBuilding(buildingId)
	if err != nil {
		return nil, err
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
		BuildingId: buildingId,
		Secret:     toInt(genNumber()),
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

func check(u *User, gameId int64, answer int) (bulls int, cows int, err error) {
	game, err := getGame(gameId)
	if err != nil {
		return
	}
	if game.UserId != u.Id {
		return 0, 0, fmt.Errorf("Forbiten")
	}

	bulls, cows, err = bullsAndCows(toByte(game.Secret), toByte(answer))
	if err != nil {
		return
	}

	build, err := getBuilding(game.BuildingId)
	if err != nil {
		return
	}

	fmt.Printf("Check game(%d): answer=%d (b:%d, c:%d)\n", game.Secret, answer, bulls, cows)

	// win
	if bulls == 4 {
		game.Status = 1
		_, err = db.Id(game.Id).Update(game)
		if err != nil {
			return
		}
		_, err = db.Exec("update building set owner_id = ? where id = ?", u.Id, game.BuildingId)
		if err != nil {
			return
		}
		_, err = db.Exec("update user set profit = profit + ? where id = ?", build.Profit, u.Id)
		if err != nil {
			return
		}
	}

	return
}
