package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zonafirmann/go-cloud-inventory/config"
	"github.com/zonafirmann/go-cloud-inventory/handlers"
)

// enableCORS adalah fungsi Middleware standar industri untuk mengizinkan akses dari Front-End
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mengizinkan semua origin (domain/port) untuk mengakses API ini
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Mengizinkan metode HTTP tertentu
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Mengizinkan header tertentu yang sering dipakai (seperti format JSON)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Menangani "Preflight Request" dari browser (pertanyaan izin awal)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Jika izin aman, lanjutkan ke fungsi handler utama
		next(w, r)
	}
}

func main() {
	// 1. Initialize Database Connection
	db := config.ConnectDB()
	defer db.Close(context.Background())

	// 2. Define API Routes (Sekarang dibungkus dengan Middleware enableCORS)
	http.HandleFunc("/products", enableCORS(handlers.GetProductsHandler(db)))
	http.HandleFunc("/checkout", enableCORS(handlers.CheckoutHandler(db)))
	http.HandleFunc("/products/analytics", enableCORS(handlers.AnalyticsHandler(db)))

	// 3. Start the Web Server
	port := ":8080"
	fmt.Printf("🚀 Cloud Inventory API Server running on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}