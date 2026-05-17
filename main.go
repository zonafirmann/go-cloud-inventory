package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5" // Menggunakan driver pgx kelas Enterprise
	"github.com/zonafirmann/go-cloud-inventory/config"
	"github.com/zonafirmann/go-cloud-inventory/handlers"
	"golang.org/x/crypto/bcrypt"
)

// Menggunakan tipe data *pgx.Conn yang benar sesuai arsitekturmu
var db *pgx.Conn
var jwtKey = []byte("zona_global_secret_key_2026")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Middleware CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	// 1. Initialize Database Connection (mengisi variabel global db)
	db = config.ConnectDB()
	// pgx membutuhkan context saat menutup koneksi
	defer db.Close(context.Background())

	// 2. Define API Routes
	http.HandleFunc("/products", enableCORS(handlers.GetProductsHandler(db)))
	http.HandleFunc("/checkout", enableCORS(authMiddleware(handlers.CheckoutHandler(db))))
	http.HandleFunc("/products/analytics", enableCORS(handlers.AnalyticsHandler(db)))
	http.HandleFunc("/register", enableCORS(registerHandler))

	// Rute Login (Dibungkus CORS agar bisa diakses React)
	http.HandleFunc("/login", enableCORS(loginHandler))

	// 3. Start the Web Server
	port := ":8080"
	fmt.Printf("🚀 Cloud Inventory API Server running on port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Jika rute ini tertangkap oleh preflight OPTIONS dari middleware, hentikan eksekusi ganda
	if r.Method == http.MethodOptions {
		return
	}

	// 1. Baca input username & password dari Front-End
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 2. Cari username di Database (Di pgx, argumen pertama WAJIB context.Background())
	var expectedPasswordHash, role string
	err = db.QueryRow(context.Background(), "SELECT password_hash, role FROM users WHERE username=$1", creds.Username).Scan(&expectedPasswordHash, &role)

	if err != nil {
		// Menggunakan error spesifik dari pgx
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 3. Bandingkan Password (Bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(expectedPasswordHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 4. Jika sukses, buat Token JWT (Masa aktif 24 jam)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 5. Kirim Token kembali ke React
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   tokenString,
		"message": "Login successful",
	})
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Izinkan CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	// 2. Baca input dari user
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 3. Enkripsi (Hash) password menggunakan Bcrypt!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 4. Simpan ke Database
	_, err = db.Exec(context.Background(), "INSERT INTO users (username, password_hash, role) VALUES ($1, $2, 'cashier')", creds.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Username already exists or database error", http.StatusInternalServerError)
		return
	}

	// 5. Beri respon sukses
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User registered successfully!"}`))
}

// authMiddleware adalah "Satpam" yang memeriksa KTP Digital (JWT) sebelum mengizinkan akses
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Periksa apakah pengunjung membawa KTP di tangannya (Header Authorization)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Akses Ditolak: Anda belum Login (Token tidak ditemukan)", http.StatusUnauthorized)
			return
		}

		// 2. Format KTP standar internasional adalah "Bearer [TOKEN_PANJANG]"
		// Kita harus membuang kata "Bearer " untuk membaca token aslinya
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			http.Error(w, "Akses Ditolak: Format token salah", http.StatusUnauthorized)
			return
		}

		// 3. Verifikasi keaslian Token menggunakan Kunci Rahasia (jwtKey) kita
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// 4. Jika token palsu, hasil editan hacker, atau sudah kadaluarsa (lewat 24 jam)
		if err != nil || !token.Valid {
			http.Error(w, "Akses Ditolak: Token palsu atau sudah kadaluarsa", http.StatusUnauthorized)
			return
		}

		// 5. Jika lolos semua pemeriksaan, silakan masuk ke ruangan (fungsi checkout/produk)
		next(w, r)
	}
}
