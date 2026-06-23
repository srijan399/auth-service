package db

import (
	"context"
	"fmt"
	"log"
)

func ConnectMain() {
	conn, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	DB = conn

	fmt.Println("Connected to Postgres!")
}

func RunQuery(queryString string) error {
	query := queryString

	_, err := DB.Exec(context.Background(), query)
	return err
}
