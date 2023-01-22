package utils

import (
	"math/rand"
	"time"
)

func ItemInSlice[T comparable](a T, list []T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	)

	s := make([]rune, n)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
