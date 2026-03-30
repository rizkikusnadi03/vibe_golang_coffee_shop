package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// NewRedis membuat koneksi baru ke Redis server dan memverifikasi koneksi.
func NewRedis(addr, password string) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Verifikasi koneksi benar-benar terbentuk
	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("gagal melakukan ping ke Redis: %w", err)
	}

	return client, nil
}
