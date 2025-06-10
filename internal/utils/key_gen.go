package utils

import (
	"crypto/rand"
	"fmt"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Генерация случайной строки заданной длины
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}

	return string(bytes), nil
}


func GenerateKey() (string, error) {
	part1, err := generateRandomString(4)
	if err != nil {
		return "", err
	}

	part2, err := generateRandomString(8)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", part1, part2), nil
}
