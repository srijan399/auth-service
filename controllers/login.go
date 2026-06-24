package controllers

import (
	"encoding/json"
	"fmt"
	"goauth/types/auth"
	dbqueries "goauth/utils/dbQueries"
	"goauth/utils/hashing"
	"goauth/utils/jwt"
	"io"
	"net/http"
)

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, "Error reading request body.", http.StatusBadRequest)
	}

	var user auth.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Fprintf(res, "Error encoding user body. Try again.")
		return
	}

	exists, existingUser, existingUserID, err := dbqueries.CheckExistingUser(user.Email)

	if err != nil {
		fmt.Fprintf(res, "Error checking existing user in db")
		return
	}

	if !exists {
		fmt.Fprintf(res, "User doesn't exist. Please sign up.")
		return
	}

	errorLogin := hashing.CompareHash(existingUser.PasswordHash, user.PasswordHash)

	if errorLogin != nil {
		http.Error(res, "Passwords do not match. Try again Later.", http.StatusBadRequest)
		return
	}

	token, err := jwt.GenerateToken(existingUserID, existingUser.Email)
	if err != nil {
		fmt.Fprintf(res, "Error generating JWT Token")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(auth.LoginResponse{
		UserObject: user,
		Token:      token,
	})
}
