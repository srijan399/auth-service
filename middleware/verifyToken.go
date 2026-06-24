package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")
var jwtSecret = []byte(SECRET_KEY)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		tokenString := req.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(res, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method/algorithm")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			fmt.Println(err)
			http.Error(res, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func getSecret() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}
