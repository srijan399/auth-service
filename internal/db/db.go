package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect() (*pgx.Conn, error) {
	connUrl := os.Getenv("DOCKER_POSTGRES_CONNECTION_URL")

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		return nil, err
	}

	fmt.Println("Set up DB Connection!")

	return conn, nil
}
