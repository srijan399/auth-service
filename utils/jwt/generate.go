package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SECRET_KEY := os.Getenv("SECRET_KEY")
	secret := []byte(SECRET_KEY)

	fmt.Println("SIGN SECRET:", os.Getenv("SECRET_KEY"))

	return token.SignedString(secret)
}
