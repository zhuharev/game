package models

import (
	"time"

	"github.com/zhuharev/game/modules/tile38"
)

type Building struct {
	Id int64 `json:"id"`

	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`

	OwnerId int64 `json:"owner_id"`

	Armor  int64 `json:"armor"`
	Profit int64 `json:"profit"`

	Updated time.Time `xorm:"updated"`
}

func createBuilding(b *Building) error {
	_, err := db.Insert(b)
	if err != nil {
		return err
	}
	return nil
}

func getBuilding(id int64) (*Building, error) {
	b := new(Building)
	has, err := db.Id(id).Get(b)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return b, nil
}

func FindBuildings(ids []int64) ([]Building, error) {
	var buildings = []Building{}
	err := db.In("id", ids).Find(&buildings)
	return buildings, err
}

func Nearby(lat, long float64) ([]Building, error) {
	m, e := tile38.Nearby(lat, long)
	if e != nil {
		return nil, e
	}
	return makeBuildingsFromMap(m), nil
}

func makeBuildingsFromMap(in map[int64][]float64) []Building {
	var res []Building
	for id, points := range in {
		if len(points) != 2 {
			continue
		}
		res = append(res, Building{
			Id:   id,
			Lat:  points[0],
			Long: points[1],
		})
	}
	return res
}
