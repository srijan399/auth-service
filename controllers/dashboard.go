package controllers

import (
	"fmt"
	"net/http"
)

func HandleDashboard(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to dashboard page. JWT verified")
}
