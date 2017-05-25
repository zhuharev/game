package fcm

import (
	"gopkg.in/maddevsio/fcm.v1"

	"github.com/zhuharev/game/modules/setting"
)

func Send(target, text string) error {
	data := map[string]string{
		"text": text,
	}
	c := fcm.NewFCM(setting.App.Fcm.Key)
	_, err := c.Send(&fcm.Message{
		Data:             data,
		RegistrationIDs:  []string{target},
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: &fcm.Notification{
			Title: "Juctvalk notify",
			Body:  text,
		},
	})
	if err != nil {
		return err
	}
	// fmt.Println("Status Code   :", response.StatusCode)
	// fmt.Println("Success       :", response.Success)
	// fmt.Println("Fail          :", response.Fail)
	// fmt.Println("Canonical_ids :", response.CanonicalIDs)
	// fmt.Println("Topic MsgId   :", response.MsgID)
	return nil
}
