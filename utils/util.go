package util

import (
	"math/rand"
	"strings"
)

var charset0 = "abcdefghijklmnopqrstuvwxyz"

// generateRandomString создает строку случайной длины из букв и символов.
func GenerateRandomString(length int) string {
	//	const charset1 = "abcdefghijklmnopqrstuvwxyz"
	//var charset = shuffleString()
	var sb strings.Builder // Исправлено здесь
	for i := 0; i < length; i++ {
		sb.WriteByte(charset0[rand.Intn(len(charset0))])
	}
	return sb.String()
}
