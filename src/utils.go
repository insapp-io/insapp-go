package main

import (
	"math/rand"
	"time"
)

func GeneratePassword() string {
	res := ""
	for i := 0; i < 2; i++ {
		res += RandomString(4)
		res += "-"
	}
	res += RandomString(4)
	return res
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLKMNOPQRSTUVWXYZ"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
