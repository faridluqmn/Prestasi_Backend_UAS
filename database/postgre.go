package database

import (
	"database/sql"
	"fmt"
	"log"

	"prestasi_backend/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectPostgre() (*sql.DB, error) {
	dsn := config.Get("POSTGRES_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping PostgreSQL: %w", err)
	}

	log.Println("✅ PostgreSQL berhasil terkoneksi")
	return db, nil
}