package main

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/kataras/iris.v6"
	"time"
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

	Balance int64 `json:"balance"`
	Goods   Goods `xorm:"-" json:"goods,omitempty"`

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
	return nil, nil
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

func handleAuth(c *iris.Context) {
	var aform AuthForm
	err := c.ReadForm(&aform)
	if err != nil {
		color.Red("%s", err)
		c.JSON(200, err.Error())
		return
	}

	u, err := AuthUser(aform)
	if err != nil {
		color.Red("%s", err)
		c.JSON(200, err.Error())
		return
	}

	c.JSON(200, struct {
		Id    int64  `json:"user_id"`
		Token string `json:"token"`
	}{
		u.Id,
		u.Token,
	})
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

func me(c *iris.Context) {
	token := c.FormValue("token")
	if token == "" {
		c.JSON(200, "token is nil")
	}
	u, err := getUserByToken(token)
	if err != nil {
		c.JSON(200, err.Error())
		return
	}

	buildings, err := getUserBuildings(u.Id)
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
	u.Goods.Buildings = buildings
	u.Goods.Count = len(buildings)

	u.Token = token
	c.JSON(200, u)
}

func getUserBuildings(useId int64) ([]Building, error) {
	var bs []Building
	err := db.Where("owner_id = ?", useId).Find(&bs)
	return bs, err
}

func getUserFromCtx(ctx *iris.Context) (*User, error) {
	token := ctx.FormValue("token")
	if token == "" {
		return nil, fmt.Errorf("Token is nil")
	}
	fmt.Println("Get user by token", token)
	return getUserByToken(token)
}

func setLocation(userId int64, lat, lon float64) error {
	_, err := db.Exec("update user set lat = ?, lon = ? where id = ?", lat, lon, userId)
	if err != nil {
		return err
	}
	return nil
}
