package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/zonafirmann/go-cloud-inventory/models"
)

// GetAllProducts mengambil semua daftar barang dari database
func GetAllProducts(conn *pgx.Conn) ([]models.Product, error) {
	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT id, name, stock, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
