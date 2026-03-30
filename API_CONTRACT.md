# API Contract — CoffeePOS

## 1. Base URL dan Versioning
Semua request dialamatkan ke Base URL berikut:  
**`http://<host>:<port>/api/v1`**

## 2. Authentication
Semua protected endpoint memerlukan pengiriman token JWT (JSON Web Token) pada header HTTP:
```http
Authorization: Bearer <jwt_token>
```
Sistem akan memberikan response `401 Unauthorized` jika token tidak ada, tidak valid, atau sudah kedaluwarsa. Response `403 Forbidden` akan diberikan jika role tidak memiliki aksess ke endpoint tersebut.

## 3. Standard Response Format

### Success Response (HTTP 200, 201)
```json
{
  "status": "success",
  "message": "Deskripsi sukses transaksi",
  "data": { ... } // Berisi object atau array sesuai entity, null jika tidak ada data
}
```
*Catatan: Endpoint dengan pagination akan membungkus array data dalam object `items` dan metadatanya dalam object `meta`.*

### Error Response (HTTP 400, 401, 403, 404, 500)
```json
{
  "status": "error",
  "message": "Deskripsi masalah yang terjadi",
  "code": "ERROR_CODE"
}
```

---

## 4. Daftar Endpoint Berdasarkan Domain

### 4.1 Auth

#### POST `/auth/login`
- **Akses:** Public
- **Deskripsi:** Login pengguna (Owner/Cashier) untuk mendapatkan JWT.
- **Request Body:**
  ```json
  {
    "email": "john@example.com",
    "password": "secretpassword"
  }
  ```
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "message": "Login berhasil",
    "data": {
      "token": "eyJhbGci...",
      "user": {
        "id": "uuid-1234",
        "name": "John Doe",
        "role": "owner"
      }
    }
  }
  ```
- **Response Error (401/404):** `INVALID_CREDENTIALS`, `USER_NOT_FOUND`, `USER_INACTIVE`

#### GET `/auth/me`
- **Akses:** Owner / Cashier
- **Deskripsi:** Mendapatkan data profil pengguna yang sedang login.
- **Response Sukses (200):** Mengembalikan data user `id`, `name`, `email`, dan `role`.

---

### 4.2 Categories (Owner)

#### GET `/categories`
- **Akses:** Owner / Cashier
- **Deskripsi:** Mengambil semua kategori produk dengan jumlah produk per kategori. (Cashier perlu ini untuk filter).
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "message": "Data kategori berhasil diambil",
    "data": [
      {
        "id": "uuid-cat-1",
        "name": "Kopi",
        "total_products": 15
      }
    ]
  }
  ```

#### POST `/categories`
- **Akses:** Owner
- **Request Body:** `{"name": "Non-Kopi"}`
- **Response Sukses (201):** Data kategori baru.
- **Response Error (400):** `VALIDATION_ERROR`, `DUPLICATE_CATEGORY_NAME`

#### PUT `/categories/:id`
- **Akses:** Owner
- **Request Body:** `{"name": "Snack & Pastry"}`
- **Response Sukses (200):** Kategori berhasil diupdate.

#### DELETE `/categories/:id`
- **Akses:** Owner
- **Response Sukses (200):** Kategori dihapus.
- **Response Error (400):** `CATEGORY_HAS_PRODUCTS` (Jika masih memiliki produk aktif terkait).

---

### 4.3 Products (Owner)

#### GET `/products`
- **Akses:** Owner / Cashier
- **Query Params:** `?category_id=uuid&search=kopi&is_active=1`
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "message": "Daftar produk",
    "data": [
      {
        "id": "uuid-prod-1",
        "name": "Americano",
        "category_id": "uuid-cat-1",
        "price": 2500000,
        "stock": 50,
        "is_active": true,
        "image_url": "https://..."
      }
    ]
  }
  ```

#### POST `/products`
- **Akses:** Owner
- **Request Format:** `multipart/form-data` (untuk upload image)
  - `name` (string)
  - `description` (string)
  - `category_id` (string)
  - `price` (int - sen)
  - `image` (file)
- **Response Sukses (201):** Data produk tersimpan.

#### PUT `/products/:id`
- **Akses:** Owner
- **Request Format:** `multipart/form-data`
- **Response Sukses (200):** Data produk diupdate.

#### DELETE `/products/:id` (Soft Delete)
- **Akses:** Owner
- **Response Sukses (200):** Produk berhasil dihapus secara logis.

#### PATCH `/products/:id/status`
- **Akses:** Owner
- **Request Body:** `{"is_active": false}`
- **Response Sukses (200):** Status aktif/nonaktif berubah.

---

### 4.4 Stock (Owner)

#### GET `/stocks/movements`
- **Akses:** Owner
- **Query Params:** `?product_id=uuid&type=IN&start_date=2026-03-01&end_date=2026-03-30`
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "message": "Riwayat stok",
    "data": [
      {
        "id": "uuid-mov-1",
        "product_name": "Americano",
        "type": "IN",
        "quantity": 20,
        "reason": "Restock beans",
        "user_name": "John Doe",
        "created_at": "2026-03-30T10:00:00Z"
      }
    ]
  }
  ```

#### POST `/stocks/adjust`
- **Akses:** Owner
- **Request Body:**
  ```json
  {
    "product_id": "uuid-prod-1",
    "type": "IN", // atau "OUT", "ADJUST"
    "quantity": 10,
    "reason": "Koreksi manual"
  }
  ```
- **Response Sukses (201):** Perubahan stok berhasil dicatat dan stok aktual terupdate.
- **Response Error (400):** `INSUFFICIENT_STOCK` (jika OUT lebih besar dari stok).

---

### 4.5 Tables (Owner)

#### GET `/tables`
- **Akses:** Owner / Cashier
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "data": [
      {
        "id": "uuid-tab-1",
        "table_number": "T1",
        "capacity": 4,
        "status": "available"
      }
    ]
  }
  ```

#### POST `/tables`
- **Akses:** Owner
- **Request Body:** `{"table_number": "T1", "capacity": 4}`
- **Response Sukses (201):** Meja berhasil ditambah.

#### PUT `/tables/:id`
- **Akses:** Owner
- **Request Body:** Mengubah number atau capacity.
- **Response Sukses (200):** Update berhasil.

#### DELETE `/tables/:id`
- **Akses:** Owner
- **Response Error (400):** `TABLE_IN_USE` (jika meja ada transaksi aktif).

---

### 4.6 Users/Cashier Management (Owner)

#### GET `/users`
- **Akses:** Owner
- **Query Params:** `?role=cashier`
- **Response Sukses (200):** Daftar pengguna sistem.

#### POST `/users`
- **Akses:** Owner
- **Request Body:** `{"name": "...", "email": "...", "password": "...", "role": "cashier"}`
- **Response Sukses (201):** User dibuat.

#### PUT `/users/:id`
- **Akses:** Owner
- **Request Body:** `{"name": "...", "email": "..."}`
- **Response Sukses (200):** User diupdate.

#### PATCH `/users/:id/status`
- **Akses:** Owner
- **Request Body:** `{"is_active": false}`
- **Response Sukses (200):** User diaktifkan/dinonaktifkan.
- **Response Error (400):** `CANNOT_DEACTIVATE_SELF`

---

### 4.7 Promos (Owner)

#### GET `/promos`
- **Akses:** Owner / Cashier
- **Query Params:** `?status=active` (Cashier umumnya hanya get actives)
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "data": [
      {
        "id": "uuid-promo-1",
        "name": "Friday Yay 10%",
        "discount_type": "percentage",
        "discount_value": 10,
        "min_purchase": 5000000,
        "max_discount": 2000000,
        "start_date": "2026-03-01T00:00:00Z",
        "end_date": "2026-03-31T23:59:59Z",
        "is_active": true
      }
    ]
  }
  ```

#### POST `/promos`
- **Akses:** Owner
- **Request Body:** Struktur sama dengan response sukses.
- **Response Sukses (201):** Promo berhasil ditambahkan.

#### PUT `/promos/:id`
- **Akses:** Owner
- **Response Sukses (200):** Promo berhasil diubah.

#### DELETE `/promos/:id` (Soft Delete)
- **Akses:** Owner
- **Response Sukses (200):** Promo dihapus secara logis.

---

### 4.8 Reports (Owner)

#### GET `/reports/dashboard`
- **Akses:** Owner
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "data": {
      "daily_revenue": 150000000,
      "weekly_revenue": [ { "date": "...", "total": 100000000 }, ... ],
      "monthly_revenue": [ ... ],
      "top_products": [ { "name": "Americano", "qty_sold": 150 } ],
      "cashier_performance": [ { "name": "Cashier A", "transactions": 45, "revenue": 500000000 } ]
    }
  }
  ```

#### GET `/reports/export`
- **Akses:** Owner
- **Query Params:** `?type=transactions&start_date=2026-03-01&end_date=2026-03-30`
- **Response Sukses (200):** File berformat `text/csv` yang di-download klien.

---

### 4.9 Shifts (Cashier)

#### GET `/shifts/active`
- **Akses:** Cashier
- **Deskripsi:** Mendapatkan shift yang saat ini sedang dibuka oleh user tersebut.
- **Response Sukses (200):**
  ```json
  {
    "status": "success",
    "data": {
      "id": "uuid-shift-1",
      "opening_cash": 10000000,
      "opened_at": "..."
    }
  }
  ```
- **Response Error (404):** `NO_ACTIVE_SHIFT`

#### POST `/shifts/start`
- **Akses:** Cashier
- **Request Body:** `{"opening_cash": 5000000}`
- **Response Sukses (201):** Shift dibuka.
- **Response Error (400):** `SHIFT_ALREADY_ACTIVE`

#### POST `/shifts/close`
- **Akses:** Cashier
- **Response Sukses (200):** 
  ```json
  {
    "status": "success",
    "message": "Shift ditutup pada akhir hari",
    "data": {
      "opening_cash": 5000000,
      "total_transactions": 25,
      "total_revenue": 150000000,
      "closing_cash": 155000000 // modal + revenue
    }
  }
  ```
- **Response Error (400):** `PENDING_TRANSACTIONS_EXIST`

---

### 4.10 Orders (Cashier)

#### GET `/orders`
- **Akses:** Cashier (hanya order di shift aktif) / Owner (semua order)
- **Response Sukses (200):** Daftar transaksi dengan summary subtotal, total, status `pending`, `paid`, `failed`, `cancelled`.

#### GET `/orders/:id`
- **Akses:** Cashier / Owner
- **Response Sukses (200):** Termasuk array `order_items` lengkap dengan snapshot nama & harga produk.

#### POST `/orders`
- **Akses:** Cashier
- **Deskripsi:** Membuat transaksi baru dan langsung memanggil prosedur checkout (Snap token).
- **Request Body:**
  ```json
  {
    "table_id": "uuid-tab-1",
    "promo_id": "uuid-promo-opstional",
    "items": [
      {
        "product_id": "uuid-prod-1",
        "quantity": 2
      }
    ]
  }
  ```
- **Response Sukses (201):**
  ```json
  {
    "status": "success",
    "data": {
      "order_id": "uuid-order-1",
      "order_number": "ORD-20260330-0001",
      "snap_token": "midtrans-snap-token-xyz123",
      "subtotal": 5000000,
      "discount_amount": 0,
      "total": 5000000,
      "status": "pending"
    }
  }
  ```
- **Response Error (400):** `INSUFFICIENT_STOCK`, `NO_ACTIVE_SHIFT`, `TABLE_IN_USE`

#### POST `/orders/:id/cancel`
- **Akses:** Cashier
- **Deskripsi:** Membatalkan transaksi yang statusnya `pending`.
- **Response Sukses (200):** Status order menjadi `cancelled`. Meja kembali `available`.

---

### 4.11 Payments

#### POST `/payments/midtrans/webhook`
- **Akses:** Public (Server-to-Server dari Midtrans)
- **Deskripsi:** Menerima notifikasi dari Midtrans. Memvalidasi Signature Key dan menerapkan state settlement/cancel.
- **Request Body:** (Sesuai dengan format callback Midtrans Core API)
  ```json
  {
    "transaction_time": "2026-03-30 15:30:00",
    "transaction_status": "settlement",
    "transaction_id": "uuid-midtrans",
    "status_message": "midtrans payment notification",
    "status_code": "200",
    "signature_key": "hashing-xxx",
    "order_id": "ORD-20260330-0001",
    "gross_amount": "50000.00"
  }
  ```
- **Response Sukses (200):** `{"status": "ok"}`
  - *Behind the scenes*: Sistem mengupdate order menjadi `paid`, memotong `stock` produk, menyimpan log di `stock_movements`, dan mengubah status meja kembali ke `available`.
- **Response Error (400/403):** `INVALID_SIGNATURE` (403), `ORDER_NOT_FOUND` (404)
