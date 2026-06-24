package auth

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}
