package models

import "fmt"

var pricesCache *Prices

// Prices represent current prices to goods
// easyjson:json
type Prices struct {
	// блокировка взлома
	Block int64
	// сброс пароля
	Reset int64
	Armor int64
}

func (p Prices) String() string {
	return fmt.Sprintf("Блокировка взлома: %d\nСброс пароля: %d\nПовышение или понижение защиты: %d", p.Block, p.Reset, p.Armor)
}

func GetPrices() (*Prices, error) {
	if pricesCache == nil {
		tmpCache := new(Prices)
		kv, err := Get("prices")
		if err != nil && err != ErrNotFound {
			return nil, err
		}
		if kv == nil {
			kv = new(Kv)
		}
		err = kv.MapTo(tmpCache)
		if err != nil {
			return nil, err
		}
		pricesCache = tmpCache
	}
	return pricesCache, nil
}

// SetPrice
func SetPrice(name string, value int64) error {
	if value <= 0 {
		return fmt.Errorf("Price zero")
	}
	switch name {
	case "block":
		pricesCache.Block = value
	case "reset":
		pricesCache.Reset = value
	case "armor":
		pricesCache.Armor = value
	}
	err := SetStruct("prices", pricesCache)
	if err != nil {
		return err
	}
	return nil
}
