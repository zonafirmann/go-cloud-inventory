package models

import "time"

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}

type Transaction struct {
	ID           int       `json:"id"`
	ProductID    int       `json:"product_id"`
	Quantity     int       `json:"quantity"`
	TotalPrice   int       `json:"total_price"`
	Status       string    `json:"status"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type InventoryReport struct {
	ProductID      int    `json:"product_id"`
	ProductName    string `json:"product_name"`
	CurrentStock   int    `json:"current_stock"`
	StatusAlert    string `json:"status_alert"`
	Recommendation string `json:"recommendation"`
}
