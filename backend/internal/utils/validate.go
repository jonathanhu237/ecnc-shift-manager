package utils

import "net/mail"

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidRole(role string) bool {
	return role == "普通助理" || role == "资深助理" || role == "黑心"
}
