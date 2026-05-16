package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/zonafirmann/go-cloud-inventory/models"
)

// GetSmartAnalysis menganalisis stok barang untuk memberikan rekomendasi bisnis
func GetSmartAnalysis(conn *pgx.Conn) ([]models.InventoryReport, error) {
	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT id, name, stock, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.InventoryReport
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock, &p.Price); err != nil {
			return nil, err
		}

		// Aturan Bisnis Bawaan
		status := "NORMAL"
		recommendation := "Stock level is healthy. No action required."

		if p.Stock <= 3 {
			status = "CRITICAL ALERT"
			recommendation = "Stock is dangerously low! Reorder immediately from supplier."
		} else if p.Stock > 50 && p.Price > 1000000 {
			status = "OVERSTOCK WARNING"
			recommendation = "Capital tied up in expensive inventory. Consider a promotional discount."
		}

		reports = append(reports, models.InventoryReport{
			ProductID:      p.ID,
			ProductName:    p.Name,
			CurrentStock:   p.Stock,
			StatusAlert:    status,
			Recommendation: recommendation,
		})
	}
	return reports, nil
}
