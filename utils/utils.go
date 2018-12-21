package utils

import (
	"log"
	"math/rand"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var chars = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func GetRandId(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
