package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/zonafirmann/go-cloud-inventory/models/repository"
)

type CheckoutReq struct {
	ProductID    int    `json:"product_id"`
	Quantity     int    `json:"quantity"`
	CustomerName string `json:"customer_name"`
}

func GetProductsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := repository.GetAllProducts(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func CheckoutHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}
		var req CheckoutReq
		json.NewDecoder(r.Body).Decode(&req)

		err := repository.CreateTransaction(db, req.ProductID, req.Quantity, req.CustomerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Transaction successful!"})
	}
}

func AnalyticsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		report, _ := repository.GetSmartAnalysis(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	}
}
