package main

import (
	"testing"
)

func TestInc(t *testing.T) {
	init := 10000

	b := Inc(Balance(init), 1)
	for i := 0; i < 120; i++ {
		b = Inc(b, 1)
	}

	if b != 120 {
		t.Errorf("b != 120, got %d", b)
	}
}
