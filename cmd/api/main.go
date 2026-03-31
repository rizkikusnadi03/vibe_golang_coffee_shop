package main

import (
	"fmt"
	"log"
	"os"

	"belajar_golang/config"
	"belajar_golang/internal/handler"
	"belajar_golang/pkg/database"
	"belajar_golang/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Inisialisasi koneksi MySQL
	db, err := database.NewMySQL(cfg.MysqlDSN())
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Inisialisasi koneksi Redis
	rdb, err := redis.NewRedis(cfg.RedisAddr(), cfg.RedisPassword)
	if err != nil {
		log.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	defer rdb.Close()

	fmt.Println("MySQL connected.")
	fmt.Println("Redis connected.")
	fmt.Println("Starting server on port :" + cfg.AppPort)

	// Inisialisasi router dan jalankan server
	r := handler.NewRouter(cfg.AppEnv)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
