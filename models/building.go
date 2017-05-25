package models

import (
	"log"
	"time"

	"github.com/zhuharev/game/modules/fixdb"
	"github.com/zhuharev/game/modules/nearbydb"
)

// Building represent building object
// Each building has area of surface
// When capturing building, user receive one-off profit. Then building has level 1
//
// User can pump level in several directions:
// - profit (K, K-lvl)
// - armor (S-lvl)
// - refresh (U-lvl) - скорость изменения всего пин-кода, раз в час он изменяется по умолчанию
// - время (T, T-lvl) - самое дорогое улучшение - сокращает время получения прибыли.
// - Итоговая формула прибыли будет: P*K*$*(t*T)
// для упрощения взлома других домов, мгновенного сброса пароля на новый, защита
// от взлома на время и возможность создавать сеть из своих домов
// (покупка соединителей, действуют только на определенном расстоянии ~ до 100м);
// - После того, как пользователь соединит N ~ 4-5 домов сетью в фигуру, все здания
// попавшие под эту фигуру переходят в его владение в исходном состоянии(первый уровень всех улучшений)
// - Данную сеть можно разрушить только захватом домов участвующих в создании границ фигуры.
// - Когда пользователь приходит захватывать здание прокачанное в защите на
// несколько блоков пинкода, то он может выбирать какой из блоков решать в данным момент;
// Это создает элемент коллективизации путем разгадывания разных блоков единомоментно несколькими игроками.
type Building struct {
	Id int64 `json:"id"`

	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`

	OwnerId   int64  `json:"owner_id"`
	OwnerName string `xorm:"-" json:"owner_name"`

	Area int `json:"area"`

	Armor      int64 `json:"armor"`
	Profit     int64 `json:"profit"`
	Refresh    int64 `json:"refresh"`
	ProfitTime int   `json:"profit_time"`

	Updated time.Time `xorm:"updated" json:"updated,omitempty"`
}

func createBuilding(b *Building) error {
	_, err := db.Insert(b)
	if err != nil {
		return err
	}
	return nil
}

func buildGetFromSQL(id int64) (*Building, error) {
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

// BuildingGet check existing building in fixdb and find it in sql
func BuildingGet(id int64) (*Building, error) {
	points, err := fixdb.Get(id)
	if err != nil {
		return nil, err
	}
	b, err := buildGetFromSQL(id)
	if err != nil {
		if err == ErrNotFound {
			return &Building{
				Id:     id,
				Lat:    points[0],
				Long:   points[1],
				Profit: 1}, nil
		}
		return nil, err
	}
	return b, nil
}

func FindBuildings(ids []int64) ([]Building, error) {
	var buildings = []Building{}
	err := db.In("id", ids).Find(&buildings)
	return buildings, err
}

// BuildingsGetByOwners return all buildings by owner ids
func BuildingsGetByOwners(ownerIds []int64) (buildings []Building, err error) {
	err = db.In("owner_id", ownerIds).Find(&buildings)
	return
}

// Nearby return buildings by passed location
func Nearby(lat, long float64) ([]Building, error) {
	m, e := nearbydb.Nearby(lat, long)
	if e != nil {
		return nil, e
	}
	for id := range m {
		_, err := fixdb.Get(id)
		if err != nil {
			log.Println("[err]", id)
		}
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

// BuildingWithOwnerCount returns total building with owners
func BuildingWithOwnerCount() (int64, error) {
	b := new(Building)
	return db.Where("owner_id != 0").Count(b)
}
