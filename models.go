package main

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *xorm.Engine
)

const (
	dbName = "db.sqlite"
)

func SetDb() {
	var err error
	db, err = xorm.NewEngine("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	db.ShowSQL(false)

	err = db.Sync2(new(User), new(Token))
	if err != nil {
		panic(err)
	}
}
