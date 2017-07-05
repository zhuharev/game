package models

import (
	"encoding/json"
)

type Balance int64

func (b *Balance) MarshalJSON() ([]byte, error) {
	whole := int64(int64(*b) / 100)
	return json.Marshal(whole)
}

func (b Balance) Real() int64 {
	return int64(int64(b) / 100)
}

func (b Balance) Add(sum int64) Balance {
	return b + Balance(sum*100)
}

func Inc(b Balance, delta int64) Balance {
	balance := int64(b)
	//var s int64 = balance / 100
	//div := balance % 100

	//div += delta

	//f := div / 60
	//o := div % 60

	return Balance(((balance/100)+((balance%100)+delta)/60)*100 + ((balance%100)+delta)%60) //(s+f)*100 + o
}

func IncreaseBalance() error {
	_, err := db.Exec("update user set balance = balance + (profit*100*15/1440)")
	//db.Iterate(bean, fun)
	return err
}

func DecreaseBalance(userID, sum int64) error {
	_, err := db.Exec("update user set balance = balance - ? where id = ?", sum*100, userID)
	//db.Iterate(bean, fun)
	return err
}
