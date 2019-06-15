package common

import (
	"math/rand"
)

const (
	// RoleAdmin is the string constant for identifying a user with admin privileges.
	RoleAdmin = "ADMIN"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString return a string of the form '[a-zA-Z0-9]{n}'.
func RandomString(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
