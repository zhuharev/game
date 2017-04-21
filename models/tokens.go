package models

import (
	"crypto/rand"
	"fmt"
	"time"
)

type Token struct {
	Token  string
	UserId int64 `xorm:"index"`

	Created time.Time `xorm:"created"`
}

func NewToken(userId int64) (*Token, error) {
	t := &Token{
		Token:  fmt.Sprintf("%x", safeRandom(16)),
		UserId: userId,
	}
	_, err := db.Insert(t)
	return t, err
}

func GetUserByToken(tok string) (*User, error) {
	var u = new(User)
	if has, err := db.Sql("select * from user where id = ( select user_id from token where token = ? )", tok).Get(u); has { // .Where("token = ?", tok).Get(&t); has {
		return u, nil
	} else if err != nil {
		return nil, err
	} else {
		return nil, ErrNotFound
	}
}

func safeRandom(l int) []byte {
	var dest = make([]byte, l)
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
	return dest
}
