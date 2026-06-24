package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"goauth/internal/db"
	"goauth/types/auth"
	dbqueries "goauth/utils/dbQueries"
	"goauth/utils/hashing"
	"goauth/utils/validation"
	"io"
	"net/http"
)

func HandleRegister(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(res, "Error reading request body. Try again.")
		return
	}

	var user auth.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Fprintf(res, "Error encoding user body. Try again.")
		return
	}

	// Check if new user
	exists, _, _, _ := dbqueries.CheckExistingUser(user.Email)
	fmt.Println("Does user exist?", exists)

	if exists {
		http.Error(res, "User already exists in DB.", http.StatusBadRequest)
		return
	}

	err = validation.ValidateEmail(user.Email)
	if err != nil {
		http.Error(res, "Not a valid email.", http.StatusBadRequest)
		return
	}

	// Convert password to password hash
	if user.PasswordHash != "" {
		hash := hashing.GetPasswordHash(user.PasswordHash)
		user.PasswordHash = string(hash)
	}

	fmt.Fprintf(res, "User registration request: %v\n", user)

	userId := insertUser(user, res)
	role := getUserRole(userId, res)

	res.WriteHeader(http.StatusCreated)
	fmt.Fprintf(res, "New User's User ID: %v\nRole: %v\n", userId, role)
}

func insertUser(user auth.User, res http.ResponseWriter) int64 {
	const insertUserQuery = `
		INSERT into users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var userId int64

	err := db.DB.QueryRow(context.Background(), insertUserQuery, user.Email, user.PasswordHash).Scan(&userId)

	if err != nil {
		http.Error(res, "Error inserting user to db. %v\n", http.StatusInternalServerError)
	}

	assignRoleQuery := `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id
		FROM roles
		WHERE name = 'user'
		ON CONFLICT DO NOTHING
	`

	_, err = db.DB.Exec(
		context.Background(),
		assignRoleQuery,
		userId,
	)

	if err != nil {
		http.Error(res, "Error in assigning roles to user. %v\n", http.StatusInternalServerError)
	}

	return userId
}

func getUserRole(userId int64, res http.ResponseWriter) string {
	var role string

	err := db.DB.QueryRow(
		context.Background(),
		`
			SELECT r.name
			FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1
	    `,
		userId,
	).Scan(&role)

	if err != nil {
		fmt.Fprintf(res, "Error in getting roles to user. %v\n", err)
	}

	return role
}
