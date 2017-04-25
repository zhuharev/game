package vk

import (
	"fmt"

	"github.com/zhuharev/vkutil"
)

// CheckToken return user id of vk.com user
func CheckToken(token string) (int, error) {
	u := vkutil.New()
	u.VkApi.AccessToken = token
	res, err := u.UsersGet(nil)
	if err != nil {
		return 0, err
	}
	if len(res) != 1 {
		return 0, fmt.Errorf("Token invalid")
	}
	return res[0].Id, nil
}
