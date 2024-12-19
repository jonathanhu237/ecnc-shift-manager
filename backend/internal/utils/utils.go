package utils

import "math/rand"

func GenerateRandomPassword(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*")
	random_password := make([]rune, length)
	for i := range random_password {
		random_password[i] = letters[rand.Intn(len(letters))]
	}
	return string(random_password)
}
