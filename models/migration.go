package models

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

var (
	currentVersion = 1
)

// Version use for migration
type Version struct {
	Id    int64
	Value int
}

// Migrate update database data
func Migrate(db *xorm.Engine) error {
	if err := db.Sync(new(Version)); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	ver := new(Version)
	has, err := db.Id(1).Get(ver)
	if err != nil {
		return err
	}
	if !has {
		ver.Value = currentVersion
		_, err := db.Insert(ver)
		if err != nil {
			return err
		}
	}
	return nil
}
