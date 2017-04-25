package models

import (
	"fmt"
	"time"

	"github.com/zhuharev/game/modules/middleware"
)

type User struct {
	Id      int64  `json:"id"`
	VkId    int64  `json:"vk_id"`
	VkToken string `json:"vk_token"`

	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`

	Username string `json:"username,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Sex      int    `json:"sex,omitempty"`

	Balance Balance `json:"balance"`
	Profit  int64   `json:"profit"`
	Goods   Goods   `xorm:"-" json:"goods,omitempty"`

	Token string `json:"token"`

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated" json:"-"`
	Deleted time.Time `xorm:"deleted" json:"-"`
}

type Goods struct {
	Count     int        `json:"count"`
	Buildings []Building `json:"buildings"`
}

type AuthForm struct {
	VkId    int64  `form:"vk_id"`
	VkToken string `form:"vk_token"`
}

// todo check token is valid
func AuthUser(af AuthForm) (*User, error) {
	var (
		u = new(User)
	)
	if af.VkId == 0 || af.VkToken == "" {
		return nil, fmt.Errorf("err data passed")
	}
	if has, err := db.Where("vk_id = ?", af.VkId).Get(u); has {
		u.VkToken = af.VkToken
		t, err := NewToken(u.Id)
		if err != nil {
			return nil, err
		}
		u.Token = t.Token
		_, err = db.Id(u.Id).Update(u)
		if err != nil {
			return nil, err
		}
		return u, nil
	} else if err != nil {
		return nil, err
	} else {
		return createUser(af)
	}
}

func createUser(af AuthForm) (*User, error) {
	u := &User{
		VkId:    af.VkId,
		VkToken: af.VkToken,
	}
	_, err := db.Insert(u)
	if err != nil {
		return nil, err
	}
	t, err := NewToken(u.Id)
	if err != nil {
		return nil, err
	}
	u.Token = t.Token
	_, err = db.Id(u.Id).Update(u)
	if err != nil {
		return nil, err
	}
	return u, err
}

type UserStore interface {
	Get()
}

func getUser(id int64) (*User, error) {
	var u = new(User)
	has, err := db.Id(id).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return u, nil
}

func GetUserBuildings(useId int64) ([]Building, error) {
	var bs []Building
	err := db.Where("owner_id = ?", useId).Find(&bs)
	return bs, err
}

func GetUserFromCtx(ctx *middleware.Context) (*User, error) {
	token := ctx.Query("token")
	if token == "" {
		return nil, fmt.Errorf("Token is nil")
	}
	return GetUserByToken(token)
}

func SetLocation(userId int64, lat, lon float64) error {
	_, err := db.Exec("update user set lat = ?, lon = ? where id = ?", lat, lon, userId)
	if err != nil {
		return err
	}
	return nil
}
