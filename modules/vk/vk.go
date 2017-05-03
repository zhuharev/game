package vk

import (
	"fmt"

	"github.com/zhuharev/vkutil"
)

// CheckToken return user id of vk.com user
func CheckToken(token string) (user vkutil.User, err error) {
	var (
		u   = vkutil.New()
		res []vkutil.User
	)
	u.VkApi.AccessToken = token
	u.VkApi.Lang = "ru"
	res, err = u.UsersGet(nil)
	if err != nil {
		return
	}
	if len(res) != 1 {
		err = fmt.Errorf("Token invalid")
		return
	}
	return res[0], nil
}
