package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQL membuat koneksi baru ke database MySQL dan mengkonfigurasi connection pool.
func NewMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi MySQL: %w", err)
	}

	// Konfigurasi connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	// Verifikasi koneksi benar-benar terbentuk
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal melakukan ping ke MySQL: %w", err)
	}

	return db, nil
}
