// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/lunny/log"
	"github.com/zhuharev/game/modules/setting"
	"github.com/zhuharev/game/modules/vk"
)

const _MIN_DB_VER = 1

type Migration interface {
	Description() string
	Migrate(*xorm.Engine) error
}

type migration struct {
	description string
	migrate     func(*xorm.Engine) error
}

func NewMigration(desc string, fn func(*xorm.Engine) error) Migration {
	return &migration{desc, fn}
}

func (m *migration) Description() string {
	return m.description
}

func (m *migration) Migrate(x *xorm.Engine) error {
	return m.migrate(x)
}

// The version table. Should have only one row with id==1
type Version struct {
	Id    int64
	Value int64
}

// This is a sequence of migrations. Add new migrations to the bottom of the list.
// If you want to "retire" a migration, remove it from the top of the list and
// update _MIN_VER_DB accordingly
var migrations = []Migration{
	// v0 -> v4 : before 0.6.0 -> last support 0.7.33
	// v4 -> v10: before 0.7.0 -> last support 0.9.141
	NewMigration("get avatar urls from vk", getAvatarUrlsFromVk),   // V10 -> V11:v0.8.5
	NewMigration("update building area from fixed.db", updateArea), // V10 -> V11:v0.8.5
}

// Migrate database to current version
func Migrate(x *xorm.Engine) error {
	if err := x.Sync(new(Version)); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	currentVersion := &Version{Id: 1}
	has, err := x.Get(currentVersion)
	if err != nil {
		return fmt.Errorf("get: %v", err)
	} else if !has {
		// If the version record does not exist we think
		// it is a fresh installation and we can skip all migrations.
		currentVersion.Id = 0
		currentVersion.Value = int64(_MIN_DB_VER + len(migrations))

		if _, err = x.InsertOne(currentVersion); err != nil {
			return fmt.Errorf("insert: %v", err)
		}
	}

	v := currentVersion.Value
	if _MIN_DB_VER > v {
		log.Fatal(0, `
Hi there, thank you for using Gogs for so long!
However, Gogs has stopped supporting auto-migration from your previously installed version.
But the good news is, it's very easy to fix this problem!
You can migrate your older database using a previous release, then you can upgrade to the newest version.

Please save following instructions to somewhere and start working:

- If you were using below 0.6.0 (e.g. 0.5.x), download last supported archive from following link:
	https://github.com/gogits/gogs/releases/tag/v0.7.33
- If you were using below 0.7.0 (e.g. 0.6.x), download last supported archive from following link:
	https://github.com/gogits/gogs/releases/tag/v0.9.141

Once finished downloading,

1. Extract the archive and to upgrade steps as usual.
2. Run it once. To verify, you should see some migration traces.
3. Once it starts web server successfully, stop it.
4. Now it's time to put back the release archive you originally intent to upgrade.
5. Enjoy!

In case you're stilling getting this notice, go through instructions again until it disappears.`)
		return nil
	}

	if int(v-_MIN_DB_VER) > len(migrations) {
		// User downgraded Gogs.
		currentVersion.Value = int64(len(migrations) + _MIN_DB_VER)
		_, err = x.Id(1).Update(currentVersion)
		return err
	}
	for i, m := range migrations[v-_MIN_DB_VER:] {
		log.Info("Migration: %s", m.Description())
		if err = m.Migrate(x); err != nil {
			return fmt.Errorf("do migrate: %v", err)
		}
		currentVersion.Value = v + int64(i) + 1
		if _, err = x.Id(1).Update(currentVersion); err != nil {
			return err
		}
	}
	return nil
}

func sessionRelease(sess *xorm.Session) {
	if !sess.IsClosed() {
		sess.Rollback()
	}
	sess.Close()
}

type User struct {
	Id        int64
	VkId      int64
	AvatarURL string `xorm:"avatar_url"`
}

func getAvatarUrlsFromVk(x *xorm.Engine) error {
	token := setting.App.Vk.AdminToken
	vkUser, err := vk.CheckToken(token)
	if err != nil {
		return err
	}
	var u User
	u.AvatarURL = vkUser.Photo200
	sess := x.Table("user").Cols("avatar_url").Where("vk_id = ?", vkUser.Id)
	_, err = sess.Update(&u)
	return err
}
