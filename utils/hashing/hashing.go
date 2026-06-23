package hashing

import (
	"fmt"
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
func CompareHash(hash string, pw string) bool {
	hashFromDB := []byte(hash)
	inputPassword := []byte(pw)
	err := bcrypt.CompareHashAndPassword(hashFromDB, inputPassword)

	if err != nil {
		fmt.Println("Error comparing hash and password. Possible mismatch.")
		return false
	}

	return true
}
