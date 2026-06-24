package dbqueries

import (
	"context"
	"database/sql"
	"goauth/internal/db"
	"goauth/types/auth"
)

func CheckExistingUser(email string) (bool, auth.User, int64, error) {
	var user auth.User
	var userId int64

	const query = `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	err := db.DB.QueryRow(context.Background(), query, email).Scan(
		&userId,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, auth.User{}, -1, nil
		}
		return false, auth.User{}, -1, err
	}

	return true, user, userId, nil
}
