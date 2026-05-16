package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zonafirmann/go-cloud-inventory/config"
	"github.com/zonafirmann/go-cloud-inventory/handlers"
)

func main() {
	db := config.ConnectDB()
	defer db.Close(context.Background())

	http.HandleFunc("/products", handlers.GetProductsHandler(db))
	http.HandleFunc("/checkout", handlers.CheckoutHandler(db))
	http.HandleFunc("/products/analytics", handlers.AnalyticsHandler(db))

	port := ":8080"
	fmt.Printf("🚀 Cloud Inventory API Server running on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
