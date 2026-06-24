package hashing

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
func GetPasswordHash(password string) []byte {
	cost, err := strconv.Atoi(os.Getenv("COST_INT"))
	if err != nil {
		log.Fatal("Error converting cost from environment Variables")
	}

	convertedPasswordInput := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(convertedPasswordInput, cost)
	if err != nil {
		log.Fatal("Error generating password hash")
	}

	return hash
}

// func CompareHashAndPassword(hashedPassword, password []byte) error
func CompareHash(hashFromDB string, rawPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(rawPassword))
	if err != nil {
		return err
	}
	return nil
}
