package db

import (
	"fmt"
	"log"
	"os"
)

func Seed() bool {
	adminPwHash := os.Getenv("ADMIN_HASH")

	const createTableQuery = `
		-- Users
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		-- Roles
		CREATE TABLE IF NOT EXISTS roles (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL
		);

		-- User ↔ Role Mapping
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id BIGINT NOT NULL,
			role_id BIGINT NOT NULL,

			PRIMARY KEY (user_id, role_id),

			CONSTRAINT fk_user
				FOREIGN KEY (user_id)
				REFERENCES users(id)
				ON DELETE CASCADE,

			CONSTRAINT fk_role
				FOREIGN KEY (role_id)
				REFERENCES roles(id)
				ON DELETE CASCADE
		);
		`

	err := RunQuery(createTableQuery)
	if err != nil {
		log.Fatalf("Error Creating initial tables. \nError: %v\n", err)
	}

	// Seeding roles table
	const seedRoleTableQuery = `
		INSERT INTO roles (name)
		VALUES ('user'), ('admin')
		ON CONFLICT (name) DO NOTHING;
	`

	err = RunQuery(seedRoleTableQuery)
	if err != nil {
		log.Fatalf("Error seeding roles table. \nError: %v\n", err)
	}

	// Inserting initial admin user
	insertAdminQuery := fmt.Sprintf(`
		INSERT INTO users (email, password_hash)
		VALUES
		('srijan@gmail.com', '%v')
		ON CONFLICT (email) DO NOTHING;
	`, adminPwHash)

	err = RunQuery(insertAdminQuery)
	if err != nil {
		log.Fatalf("Error adding admin user. \nError: %v\n", err)
	}

	// Assign admin role
	const assignAdmin = `INSERT INTO user_roles (user_id, role_id)
						SELECT u.id, r.id
						FROM users u, roles r
						WHERE u.email = 'srijan@gmail.com'
						AND r.name = 'admin'
						ON CONFLICT DO NOTHING;
						`
	err = RunQuery(assignAdmin)
	if err != nil {
		log.Fatalf("Error assigning admin user. \nError: %v\n", err)
	}

	return true
}
