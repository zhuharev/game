package models

import (
	"encoding/json"
)

// Kv represent key-value data
type Kv struct {
	Key   string `xorm:"pk"`
	Value string
}

// Set value to key-value store
func Set(kv *Kv) error {
	sql := `INSERT OR REPLACE INTO kv (key, value) VALUES ( ?, ?);`
	_, err := db.Exec(sql, kv.Key, kv.Value)
	return err
}

// Set value to key-value store
func SetStruct(key string, str interface{}) error {
	kv, err := NewKvFromStruct(key, str)
	if err != nil {
		return err
	}
	sql := `INSERT OR REPLACE INTO kv (key, value) VALUES ( ?, ?);`
	_, err = db.Exec(sql, kv.Key, kv.Value)
	return err
}

// SetString value to key-value store
func SetString(key, value string) error {
	sql := `INSERT OR REPLACE INTO kv (key, value) VALUES ( ?, ?);`
	_, err := db.Exec(sql, key, value)
	return err
}

// Get receive value from kv store
func Get(key string) (*Kv, error) {
	kv := new(Kv)
	has, err := db.Where("key = ?", key).Get(kv)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return kv, nil
}

// GetString receive value from kv store
func GetString(key string) (string, error) {
	kv := new(Kv)
	has, err := db.Where("key = ?", key).Get(kv)
	if err != nil {
		return "", err
	}
	if !has {
		return "", ErrNotFound
	}
	return kv.Value, nil
}

// MapTo map value as bytes json
func (kv *Kv) MapTo(target interface{}) error {
	if kv != nil && kv.Value == "" {
		return nil
	}
	return json.Unmarshal([]byte(kv.Value), target)
}

func NewKvFromStruct(key string, str interface{}) (*Kv, error) {
	bts, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	kv := new(Kv)
	kv.Value = string(bts)
	kv.Key = key
	return kv, nil
}
