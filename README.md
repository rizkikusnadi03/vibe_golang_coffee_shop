# Coffee Shop POS API

A robust **Point of Sales (POS)** backend service tailored for a medium-scale coffee shop, written in Go. This repository houses comprehensive database schemas, structured routing setups, custom semantic validations in Indonesian, and a standardized REST JSON response pattern.

## Features
- **User Authentication and Roles**: Distinguishes between `owner` and `cashier` with customized access levels.
- **Product Management**: Manage categories, products, stock, pricing, and statuses flexibly.
- **Order and Payment System**: Record orders by table or takeaway, calculate subtotal and discount, and handle transactions.
- **Promo Mechanism**: Handles percentage or fixed discount, usage limits, and active periods.
- **Shift & Till Management**: Cashier shifts with opening and closing cash balancing.
- **Extensive Semantic Validations**: Robust and human-readable payload validation in Indonesian language via custom error tag mapping.

## Tech Stack
- **Language**: Go 1.25+
- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database**: MySQL 8.0+
- **Migration Engine**: [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- **Validation**: [go-playground/validator v10](https://github.com/go-playground/validator)

## Getting Started

### Prerequisites
- Go installed on your machine.
- Local or remote instance of MySQL.
- GNU Make.

### Available Commands (Make)
You can use `make` commands to easily bootstrap and run the application.

```bash
# Run the development server
make run

# Clean and update Go module dependencies
make tidy
```

#### Database Migrations
By default, the migrations target the Database URL dynamically passed through `DB_URL`. The default is set to `mysql://root:secret@tcp(localhost:3306)/coffee_pos`.

```bash
# Run all database migrations to the latest version
make migrate-up

# Rollback the last migration
make migrate-down

# Rollback all migrations (WARNING: deletes all tables)
make migrate-down-all

# Check the current active migration version
make migrate-version

# Create a new blank migration file pair
make migrate-create name=migration_name_here
```

## Directory Structure
- `cmd/api/` — Application entrypoint (`main.go`).
- `internal/handler/` — HTTP layer, router initialization, and endpoint handlers.
- `migrations/` — Structured `.up.sql` and `.down.sql` files building the Coffee Shop POS database architecture.
- `pkg/response/` — Standardized HTTP JSON Response format wrappers (e.g., `response.OK()`, `response.BadRequest()`).
- `pkg/validator/` — Payload checker returning mapped error dictionaries in formatted Indonesian rules.

## Standardized JSON Response
All successful requests yield an explicit schema like below:

```json
{
    "success": true,
    "message": "Transaksi berhasil disimpan",
    "data": { ... }
}
```

Validation errors map field names to precise contexts:
```json
{
    "success": false,
    "message": "Validasi gagal",
    "errors": {
        "email": "format email tidak valid",
        "product_id": "wajib diisi",
        "quantity": "minimal 1"
    }
}
```

## Setup using Docker (Optional)
If a `docker-compose.yml` is present in your environment (for MySQL configurations):
1. Boot the environment: `docker-compose up -d`
2. Run database migrations: `make migrate-up`
3. Serve logic: `make run`
