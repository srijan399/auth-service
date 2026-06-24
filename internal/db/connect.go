package db

import (
	"context"
	"fmt"
	"log"
	"os"
)

func ConnectMain() {
	conn, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	DB = conn

	pgAdminGateway := os.Getenv("PGADMIN_GATEWAY")

	fmt.Printf("Connected to Postgres. Go to pgAdmin at %v\n", pgAdminGateway)
}

func RunQuery(queryString string) error {
	query := queryString

	_, err := DB.Exec(context.Background(), query)
	return err
}
