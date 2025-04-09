package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword - take user's password as input, generates a bcrypt hash with a cost factor of default
// VerifyPassword - take password, and a stored hash as inputs, compare password against the hash and return true/false

func HashPassword(passwordPlain string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordPlain),bcrypt.DefaultCost)
	return string(passwordHash),err
}

func VerifyPassword(passwordAttempt, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordAttempt))
	return err == nil
}