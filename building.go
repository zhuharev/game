package main

import (
	"time"
)

type Building struct {
	Id int64 `json:"id"`

	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`

	OwnerId int64 `json:"owner_id"`

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

func findBuildings(ids []int64) ([]Building, error) {
	var buildings = []Building{}
	err := db.In("id", ids).Find(&buildings)
	return buildings, err
}
