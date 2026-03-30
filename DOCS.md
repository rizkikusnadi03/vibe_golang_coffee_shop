# Dokumentasi File — Coffee Shop POS

Dokumen ini menjelaskan setiap file yang telah dibuat dalam project Coffee Shop POS, termasuk tujuan, isi, dan cara penggunaannya.

---

## Daftar Isi

1. [Dokumentasi Bisnis](#1-dokumentasi-bisnis)
2. [Entry Point Aplikasi](#2-entry-point-aplikasi)
3. [Konfigurasi](#3-konfigurasi)
4. [Package Database & Redis](#4-package-database--redis)
5. [Docker & Infrastructure](#5-docker--infrastructure)
6. [API & Testing Tools](#6-api--testing-tools)
7. [Struktur Folder Internal](#7-struktur-folder-internal)

---

## 1. Dokumentasi Bisnis

### `PRD.md`

- **Tujuan:** Product Requirements Document — mendefinisikan kebutuhan bisnis sistem POS.
- **Isi:** Deskripsi produk, user roles (Owner & Cashier), fitur lengkap per role, acceptance criteria, dan 30+ business rules.
- **Digunakan oleh:** Tim development sebagai acuan utama fitur yang harus diimplementasikan.

### `SCHEMA.md`

- **Tujuan:** Desain database MySQL lengkap.
- **Isi:** 9 tabel (`users`, `categories`, `products`, `stock_movements`, `tables`, `shifts`, `promos`, `orders`, `order_items`), relasi antar tabel (FK), indeks, dan keputusan desain teknis (UUID, BIGINT untuk uang, soft delete).
- **Digunakan oleh:** Developer backend saat membuat migration dan repository layer.

### `API_CONTRACT.md`

- **Tujuan:** Kontrak API endpoint lengkap.
- **Isi:** 11 domain endpoint (Auth, Categories, Products, Stock, Tables, Users, Promos, Reports, Shifts, Orders, Payments) dengan method, path, request/response body dalam format JSON, dan error codes.
- **Digunakan oleh:** Developer backend dan frontend sebagai panduan integrasi API.

---

## 2. Entry Point Aplikasi

### `cmd/api/main.go`

- **Package:** `main`
- **Tujuan:** Entry point aplikasi Go. Titik awal eksekusi program.
- **Fungsi:**
  - `main()` — Memanggil `config.Load()` untuk memuat konfigurasi dari environment variables. Jika gagal, mencetak pesan error dan keluar dengan exit code 1. Jika berhasil, mencetak konfirmasi port yang digunakan.
- **Import internal:** `belajar_golang/config`
- **Cara menjalankan:**
  ```bash
  go run cmd/api/main.go
  ```

---

## 3. Konfigurasi

### `config/config.go`

- **Package:** `config`
- **Tujuan:** Memuat dan memvalidasi semua environment variables yang dibutuhkan aplikasi.
- **Struct `Config`:** Menyimpan semua konfigurasi aplikasi dalam 5 grup:

  | Grup | Field | Tipe |
  |------|-------|------|
  | App | `AppPort`, `AppEnv` | `string` |
  | Database | `DBHost`, `DBPort`, `DBName`, `DBUser`, `DBPassword` | `string` |
  | Redis | `RedisHost`, `RedisPort`, `RedisPassword` | `string` |
  | JWT | `JWTSecret`, `JWTExpiryHours` | `string`, `int` |
  | Midtrans | `MidtransServerKey`, `MidtransClientKey`, `MidtransEnv` | `string` |

- **Fungsi:**
  - `Load() (*Config, error)` — Memuat `.env` file (opsional, tidak error jika tidak ada), membaca semua `os.Getenv()`, menerapkan default values (`AppPort=8080`, `DBPort=3306`, `RedisPort=6379`, `JWTExpiryHours=24`), lalu memvalidasi field wajib (`DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `JWT_SECRET`).
  - `MysqlDSN() string` — Menghasilkan connection string format: `user:password@tcp(host:port)/dbname?parseTime=true`
  - `RedisAddr() string` — Menghasilkan alamat Redis format: `host:port`
- **Dependency eksternal:** `github.com/joho/godotenv`

### `.env.example`

- **Tujuan:** Template environment variables. Developer baru meng-copy file ini menjadi `.env` lalu mengisi nilainya.
- **Isi:** 15 variabel mencakup konfigurasi App, Database, Redis, JWT, dan Midtrans.
- **Cara pakai:**
  ```bash
  cp .env.example .env
  # Edit .env dan isi nilai yang sesuai
  ```

### `.env`

- **Tujuan:** File environment variables aktif untuk development (tidak di-commit ke Git).
- **Isi:** Nilai aktual untuk semua variabel yang didefinisikan di `.env.example`, termasuk variabel tambahan `MYSQL_ROOT_PASSWORD` dan `MYSQL_DATABASE` untuk Docker Compose.

---

## 4. Package Database & Redis

### `pkg/database/mysql.go`

- **Package:** `database`
- **Tujuan:** Membuat dan mengkonfigurasi koneksi ke MySQL.
- **Fungsi:**
  - `NewMySQL(dsn string) (*sql.DB, error)` — Membuka koneksi MySQL menggunakan DSN string, mengatur connection pool, lalu memverifikasi koneksi dengan `Ping()`.
- **Connection Pool Settings:**

  | Parameter | Nilai | Penjelasan |
  |-----------|-------|------------|
  | `MaxOpenConns` | 25 | Maksimal koneksi aktif bersamaan |
  | `MaxIdleConns` | 10 | Maksimal koneksi idle di pool |
  | `ConnMaxLifetime` | 5 menit | Umur maksimal satu koneksi |
  | `ConnMaxIdleTime` | 3 menit | Waktu idle maksimal sebelum koneksi ditutup |

- **Error Handling:** Jika `Ping()` gagal, koneksi langsung di-`Close()` dan return error deskriptif.
- **Dependency:** `github.com/go-sql-driver/mysql` (blank import untuk registrasi driver)
- **Contoh penggunaan:**
  ```go
  db, err := database.NewMySQL(cfg.MysqlDSN())
  ```

### `pkg/redis/redis.go`

- **Package:** `redis`
- **Tujuan:** Membuat dan memverifikasi koneksi ke Redis server.
- **Fungsi:**
  - `NewRedis(addr, password string) (*goredis.Client, error)` — Membuat Redis client dari library `go-redis/v9` dengan DB index 0, lalu memverifikasi koneksi dengan `Ping()`.
- **Import alias:** Menggunakan `goredis` sebagai alias untuk menghindari konflik dengan nama package sendiri (`redis`).
- **Error Handling:** Jika `Ping()` gagal, client langsung di-`Close()` dan return error deskriptif.
- **Dependency:** `github.com/redis/go-redis/v9`
- **Contoh penggunaan:**
  ```go
  rdb, err := redis.NewRedis(cfg.RedisAddr(), cfg.RedisPassword)
  ```

---

## 5. Docker & Infrastructure

### `docker-compose.yml` (Development)

- **Tujuan:** Menjalankan infrastruktur development (MySQL + Redis) secara lokal.
- **Services:**

  | Service | Image | Port Host | Volume |
  |---------|-------|-----------|--------|
  | `mysql` | `mysql:8.0` | `127.0.0.1:3306` | `mysql_data` |
  | `redis` | `redis:7-alpine` | `127.0.0.1:6379` | `redis_data` |

- **Healthcheck:** MySQL menggunakan `mysqladmin ping` (interval 10s, start_period 30s). Redis menggunakan `redis-cli ping`.
- **Network:** `coffee-pos-network` (bridge driver).
- **Cara menjalankan:**
  ```bash
  docker compose up -d
  ```

### `docker-compose.prod.yml` (Production)

- **Tujuan:** Menjalankan stack production lengkap (MySQL + Redis + App Go).
- **Perbedaan dari development:**

  | Aspek | Development | Production |
  |-------|-------------|------------|
  | Port MySQL/Redis ke host | Terbuka (`127.0.0.1`) | Tidak ada (terisolasi) |
  | Service app | Tidak ada | Ada (`build: .`) |
  | Redis auth | Tanpa password | Dengan `--requirepass` |
  | App dependency | Manual | `depends_on: service_healthy` |

- **Cara menjalankan:**
  ```bash
  docker compose -f docker-compose.prod.yml up -d --build
  ```

---

## 6. API & Testing Tools

### `postman_collection.json`

- **Tujuan:** Postman Collection yang bisa langsung diimport untuk testing semua API endpoint.
- **Isi:** 11 folder (Auth, Categories, Products, Stock, Tables, Cashiers, Promos, Reports, Shifts, Orders, Payments) dengan total 25+ request yang sudah terisi method, URL, headers, dan example body.
- **Variable:** `{{base_url}}` untuk base URL, `{{token}}` untuk JWT token.

### `postman_environment_local.json`

- **Tujuan:** Environment Postman untuk development lokal.
- **Nilai:** `base_url = http://localhost:8080/api/v1`

### `postman_environment_production.json`

- **Tujuan:** Environment Postman untuk production.
- **Nilai:** `base_url = https://yourdomain.com/api/v1`

### Cara Import ke Postman

1. Buka Postman → klik **Import**
2. Pilih ketiga file JSON sekaligus
3. Pilih environment **Local** atau **Production** di dropdown kanan atas

---

## 7. Struktur Folder Internal

Folder-folder berikut sudah disiapkan dengan file `.gitkeep` agar ter-track oleh Git. Masing-masing akan diisi pada fase implementasi berikutnya:

| Folder | Tujuan | Akan Berisi |
|--------|--------|-------------|
| `internal/entity/` | Domain struct / model | Struct Go yang merepresentasikan tabel database |
| `internal/repository/` | Repository interfaces | Interface untuk abstraksi akses data |
| `internal/repository/mysql/` | Implementasi MySQL | Implementasi konkret dari repository interfaces |
| `internal/service/` | Business logic | Service layer yang mengolah logika bisnis |
| `internal/handler/` | HTTP handlers | Handler untuk menerima dan merespons HTTP request |
| `internal/middleware/` | Middleware | Auth middleware, logging, CORS, dll |
| `internal/dto/` | Data Transfer Objects | Struct untuk request body dan response body |
| `pkg/jwt/` | JWT helper | Fungsi generate dan validate JWT token |
| `pkg/response/` | Response helper | Fungsi standar untuk format JSON response |
| `pkg/validator/` | Input validator | Validasi input request |
| `migrations/` | SQL migrations | File `.sql` untuk create/alter table |

---

## File Pendukung Lainnya

### `.gitignore`

- **Tujuan:** Mencegah file sensitif dan tidak perlu ter-commit ke Git.
- **Yang di-ignore:** `.env`, `vendor/`, `*.exe`, `tmp/`, serta file dokumentasi yang sudah di-commit secara terpisah.

### `README.md`

- **Tujuan:** Halaman utama repository di GitHub.
- **Isi:** Nama project, deskripsi singkat, dan referensi ke dokumen teknis (PRD, SCHEMA, API_CONTRACT).

### `go.mod` & `go.sum`

- **Tujuan:** Go module definition dan dependency lock file.
- **Module name:** `belajar_golang`
- **Dependencies:**
  - `github.com/joho/godotenv` — Memuat file `.env`
  - `github.com/go-sql-driver/mysql` — MySQL driver
  - `github.com/redis/go-redis/v9` — Redis client
