package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	AppEnv            string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	RedisHost         string
	RedisPort         string
	RedisPassword     string
	JWTSecret         string
	JWTExpiryHours    int
	MidtransServerKey string
	MidtransClientKey string
	MidtransEnv       string
}

// Load reads configuration from .env and environment variables.
func Load() (*Config, error) {
	// Memuat file .env secara optional
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:           os.Getenv("APP_PORT"),
		AppEnv:            os.Getenv("APP_ENV"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBName:            os.Getenv("DB_NAME"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		RedisHost:         os.Getenv("REDIS_HOST"),
		RedisPort:         os.Getenv("REDIS_PORT"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		MidtransServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransClientKey: os.Getenv("MIDTRANS_CLIENT_KEY"),
		MidtransEnv:       os.Getenv("MIDTRANS_ENV"),
	}

	// Fallback nilai default untuk Ports jika belum diset
	
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "3306"
	}
	if cfg.RedisPort == "" {
		cfg.RedisPort = "6379"
	}

	// Parsing dan set opsi default JWT Expiry Hours ke 24 jam
	expiryStr := os.Getenv("JWT_EXPIRY_HOURS")
	if expiryStr == "" {
		cfg.JWTExpiryHours = 24
	} else {
		expiry, err := strconv.Atoi(expiryStr)
		if err != nil || expiry <= 0 {
			cfg.JWTExpiryHours = 24
		} else {
			cfg.JWTExpiryHours = expiry
		}
	}

	// Validasi field yang mutlak dibutuhkan
	if cfg.DBHost == "" {
		return nil, fmt.Errorf("konfigurasi wajib hilang: DB_HOST belum diatur")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("konfigurasi wajib hilang: DB_NAME belum diatur")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("konfigurasi wajib hilang: DB_USER belum diatur")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("konfigurasi wajib hilang: DB_PASSWORD belum diatur")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("konfigurasi wajib hilang: JWT_SECRET belum diatur")
	}

	return cfg, nil
}

// MysqlDSN mem-format connection string MySQL DSN.
func (c *Config) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// RedisAddr mem-format alamat Redis konvensional.
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
