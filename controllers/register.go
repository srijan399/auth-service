package controllers

import (
	"fmt"
	"net/http"
)

func HandleRegister(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling user registration")
	// To-do: Set up registration endpoint
	// Take in email, pw, hash it and write an Insert Query SQL using DB Exec, check if user already exists, if not validate email, register
}
