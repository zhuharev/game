package models

import (
	"testing"
)

func KvTest(t *testing.T) {
	SetDb()
	kv := new(Kv)
	kv.Key = "alo"
	kv.Value = "hallo"

	err := SetString(kv.Key, kv.Value)
	if err != nil {
		t.Fatal(err)
	}

	val, err := GetString(kv.Key)
	if err != nil {
		t.Fatal(err)
	}
	if val != kv.Value {
		t.Fatalf("Val (%s)", val)
	}

	// update
	kv.Value = "opa"
	err = SetString(kv.Key, kv.Value)
	if err != nil {
		t.Fatal(err)
	}

	val, err = GetString(kv.Key)
	if err != nil {
		t.Fatal(err)
	}
	if val != kv.Value {
		t.Fatalf("Val (%s)", val)
	}
}
