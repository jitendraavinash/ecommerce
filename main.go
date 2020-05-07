package main

import (
	"ecommerce/db"
	"ecommerce/router"
	"net/http"
)

func main() {
	// establish database connection
	db.Connect()

	// route handlers
	http.HandleFunc("/health", router.ServerStatus)
	http.HandleFunc("/vendor", router.Vendor)
	http.HandleFunc("/item", router.Item)

	http.ListenAndServe(":8080", nil)
}
