package models

import (
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"time"

	"github.com/zhuharev/game/modules/tile38"
)

var (
	db          *xorm.Engine
	ErrNotFound = errors.New("not found")
)

const (
	dbName        = "db.sqlite"
	allowedLength = 4
)

func SetDb() {
	var err error
	var has bool
	if com.IsExist(dbName) {
		has = true
	}

	db, err = xorm.NewEngine("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	db.ShowSQL(false)

	err = db.Sync2(
		new(User),
		new(Token),
		new(Building),
		new(Game),
	)
	if err != nil {
		panic(err)
	}

	if !has {
		tile38.GetBuildings()
	}
}

func toInt(in []byte) int {
	return com.StrTo(string(in)).MustInt()
}

func toByte(in int) []byte {
	return []byte(fmt.Sprint(in))
}

func genNumber() []byte {
	pat := make([]byte, allowedLength)
	rand.Seed(time.Now().Unix())
	r := rand.Perm(9)
	offset := 0

	for r[0] == 0 { // yes, kind of hacky, no guarantees for time complexity here
		r = rand.Perm(9)
	}

	for i := range pat {
		pat[i] = '0' + byte(r[i+offset])
	}

	return pat
}
