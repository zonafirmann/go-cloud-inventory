package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// CreateTransaction menangani proses checkout dengan pengamanan transaksi database
func CreateTransaction(conn *pgx.Conn, productID int, qty int, customer string) error {
	ctx := context.Background()

	// Mulai Transaksi (Begin)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	// Jika ada error di tengah jalan, batalkan semua perubahan (Rollback)
	defer tx.Rollback(ctx)

	// Kunci baris data (FOR UPDATE) agar tidak ada yang bisa beli barang yang sama di waktu bersamaan
	var stock, price int
	err = tx.QueryRow(ctx, "SELECT stock, price FROM products WHERE id = $1 FOR UPDATE", productID).Scan(&stock, &price)
	if err != nil {
		return errors.New("product not found")
	}

	// Validasi Stok
	if stock < qty {
		return errors.New("insufficient stock")
	}

	// Kurangi Stok
	newStock := stock - qty
	_, err = tx.Exec(ctx, "UPDATE products SET stock = $1 WHERE id = $2", newStock, productID)
	if err != nil {
		return err
	}

	// Catat Transaksi
	totalPrice := price * qty
	_, err = tx.Exec(ctx, "INSERT INTO transactions (product_id, quantity, total_price, customer_name, status) VALUES ($1, $2, $3, $4, $5)", productID, qty, totalPrice, customer, "success")
	if err != nil {
		return err
	}

	// Simpan semua perubahan secara permanen (Commit)
	return tx.Commit(ctx)
}
