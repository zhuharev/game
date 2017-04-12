package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"math/rand"
	"time"
)

var alphabet = []byte("0123456789")

const allowedLength = 4

func validLength(guess []byte) bool {
	return len(guess) != allowedLength
}

func beginWithZero(guess []byte) bool {
	return len(guess) > 0 && guess[0] == '0'
}

func isValid(c byte) bool {
	return bytes.IndexByte(alphabet, c) != -1
}

func bullsAndCows(target []byte, guess []byte) (int, int, error) {
	var bulls, cows int

	if validLength(guess) {
		return 0, 0, fmt.Errorf("Guess must be %d characters long", allowedLength)
	}

	if beginWithZero(guess) {
		return 0, 0, fmt.Errorf("Number can not start with zero")
	}

	for idx, c := range guess {

		if !isValid(c) {
			return 0, 0, errors.New("Your guess has invalid characters")
		}

		if bytes.IndexByte(guess[:idx], c) >= 0 {
			return 0, 0, fmt.Errorf("Repeated '%c'. No repetition allowed", c)
		}

		pos := bytes.IndexByte(target, c)

		if pos == idx {
			bulls = bulls + 1
		} else {
			if pos >= 0 {
				cows = cows + 1
			}
		}
	}
	return bulls, cows, nil
}

func genNumber() []byte {
	pat := make([]byte, allowedLength)
	rand.Seed(time.Now().Unix())
	r := rand.Perm(9)
	offset := 0

	for r[0] == 0 { // yes, kind of hacky, no guarantees for time complexity here
		r = rand.Perm(9)
	}

	for i := range pat {
		pat[i] = '0' + byte(r[i+offset])
	}

	return pat
}

func toInt(in []byte) int {
	return com.StrTo(string(in)).MustInt()
}

func toByte(in int) []byte {
	return []byte(fmt.Sprint(in))
}
