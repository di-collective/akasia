package utils

import "math/rand"

const (
	alphabeticalLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumericLetters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandAlphabeticalString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabeticalLetters[rand.Intn(len(alphabeticalLetters))]
	}
	return string(b)
}

func RandAlphanumericString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanumericLetters[rand.Intn(len(alphanumericLetters))]
	}
	return string(b)
}
