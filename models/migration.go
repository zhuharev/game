package models

import (
	"github.com/go-xorm/xorm"
	"github.com/zhuharev/game/models/migrations"
)

// Migrate update database data
func Migrate(db *xorm.Engine) error {
	return migrations.Migrate(db)
}
