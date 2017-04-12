package main

import (
	"errors"
	"github.com/Unknwon/com"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db          *xorm.Engine
	ErrNotFound = errors.New("not found")
)

const (
	dbName = "db.sqlite"
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
		getBuildings()
	}
}
