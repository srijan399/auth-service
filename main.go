package main

import (
	"fmt"
	"goauth/internal/db"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Listening on port :8090")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading Environment Variables.")
		log.Fatal()
	}

	// Connect Postgres DB
	db.ConnectMain()

	// Seed Database with data
	res := db.Seed()
	if res {
		fmt.Printf("========== Seed successful ==========\n\n")
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hey, welcome to Auth Service 101")
	})

	http.ListenAndServe(":8090", nil)
}
