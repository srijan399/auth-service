package main

import (
	"fmt"
	"goauth/controllers"
	"goauth/internal/db"
	"goauth/middleware"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Hey, welcome to Auth Service 101")
	})

	http.HandleFunc("/auth/register", controllers.HandleRegister)
	http.HandleFunc("/auth/login", controllers.HandleLogin)

	protected := http.NewServeMux()
	protected.HandleFunc("/dashboard", controllers.HandleDashboard)

	http.Handle("/protected/", http.StripPrefix("/protected", middleware.AuthMiddleware(protected)))
	http.ListenAndServe(":8090", nil)
}
