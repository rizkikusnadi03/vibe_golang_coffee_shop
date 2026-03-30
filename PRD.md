# Product Requirements Document (PRD)
# Sistem Point of Sales — Coffee Shop

| Field            | Detail                                      |
|------------------|---------------------------------------------|
| **Nama Produk**  | CoffeePOS                                   |
| **Versi**        | 1.0                                         |
| **Tech Stack**   | Go (Golang), MySQL, Midtrans Snap           |
| **Tanggal**      | 30 Maret 2026                               |
| **Status**       | Draft                                       |

---

## 1. Deskripsi Produk

CoffeePOS adalah sistem Point of Sales (POS) berbasis web yang dirancang khusus untuk coffee shop skala menengah. Sistem ini membantu pemilik usaha (Owner) dalam mengelola operasional harian — mulai dari manajemen produk, stok, meja, hingga pelaporan keuangan — sekaligus menyediakan antarmuka kasir (Cashier) yang efisien untuk memproses transaksi pelanggan secara cepat dan akurat.

### Target Pengguna

| Pengguna         | Deskripsi                                                                                  |
|------------------|--------------------------------------------------------------------------------------------|
| **Owner**        | Pemilik atau manajer coffee shop yang bertanggung jawab atas operasional dan strategi bisnis |
| **Cashier**      | Staff kasir yang bertugas melayani transaksi pelanggan di outlet                             |

### Tujuan Utama

1. Menyederhanakan proses transaksi penjualan di coffee shop
2. Memberikan Owner visibilitas penuh atas operasional dan keuangan bisnis
3. Mengotomasi pencatatan stok, transaksi, dan pelaporan
4. Mendukung pembayaran digital melalui payment gateway Midtrans

---

## 2. User Roles

### 2.1 Owner

Owner memiliki **akses penuh** ke seluruh fitur sistem. Role ini ditujukan untuk pemilik atau manajer yang bertanggung jawab atas strategi dan operasional bisnis.

| Akses                      | Deskripsi                                                            |
|----------------------------|----------------------------------------------------------------------|
| Dashboard & Laporan        | Melihat ringkasan revenue, produk terlaris, dan performa cashier     |
| Manajemen Produk           | CRUD produk, upload foto, atur status aktif/nonaktif, soft delete    |
| Manajemen Kategori         | CRUD kategori produk                                                 |
| Manajemen Stok             | Penyesuaian stok manual, riwayat pergerakan stok                     |
| Manajemen Meja             | CRUD meja, atur status ketersediaan                                  |
| Manajemen User (Cashier)   | CRUD akun cashier, aktivasi/deaktivasi akun                          |
| Manajemen Promo & Diskon   | CRUD promo, atur periode aktif, tipe diskon (persentase/nominal)     |
| Export Laporan             | Download laporan dalam format CSV                                    |

### 2.2 Cashier

Cashier memiliki **akses terbatas** yang berfokus pada operasional transaksi harian.

| Akses                       | Deskripsi                                                          |
|-----------------------------|--------------------------------------------------------------------|
| Manajemen Shift             | Buka dan tutup shift harian dengan pencatatan modal kas             |
| Buat Transaksi              | Pilih meja, tambah produk ke keranjang, atur quantity               |
| Apply Promo                 | Menerapkan promo/diskon aktif ke transaksi                          |
| Checkout (Midtrans)         | Proses pembayaran melalui Midtrans Snap                            |
| Riwayat Transaksi Shift     | Melihat daftar transaksi dalam shift yang sedang berjalan hari ini  |

> [!IMPORTANT]
> Cashier **tidak dapat** mengakses fitur manajemen (produk, kategori, stok, meja, user, promo) maupun dashboard laporan.

---

## 3. Fitur Owner

### 3.1 Manajemen Produk

Owner dapat mengelola seluruh produk yang dijual di coffee shop.

#### Fungsionalitas

| Fitur               | Deskripsi                                                                                |
|----------------------|------------------------------------------------------------------------------------------|
| **Create**           | Tambah produk baru dengan nama, deskripsi, harga, kategori, foto, dan status awal        |
| **Read**             | Lihat daftar produk dengan filter (kategori, status) dan pencarian berdasarkan nama       |
| **Update**           | Edit informasi produk termasuk nama, harga, deskripsi, kategori, dan foto                 |
| **Soft Delete**      | Hapus produk secara logis (set `deleted_at`), produk tidak tampil di kasir tapi data tetap tersimpan |
| **Upload Foto**      | Upload foto produk (format: JPG, PNG; maks: 2 MB)                                        |
| **Status Aktif/Nonaktif** | Toggle status produk; produk nonaktif tidak muncul di daftar kasir                   |

#### Acceptance Criteria

- [ ] Owner dapat menambah produk baru dengan semua field wajib (nama, harga, kategori) dan produk tersimpan di database
- [ ] Owner dapat melihat daftar produk dengan pagination, filter berdasarkan kategori dan status, serta pencarian berdasarkan nama
- [ ] Owner dapat mengedit produk yang sudah ada dan perubahan langsung tercermin di sistem
- [ ] Soft delete mengisi kolom `deleted_at` dengan timestamp; produk yang di-soft-delete tidak muncul di daftar kasir maupun daftar produk Owner (kecuali ada filter khusus)
- [ ] Foto produk berhasil di-upload dan ditampilkan; sistem menolak file selain JPG/PNG atau berukuran lebih dari 2 MB
- [ ] Produk dengan status "nonaktif" tidak muncul di daftar produk kasir saat membuat transaksi

---

### 3.2 Manajemen Kategori Produk

Owner dapat mengelola kategori untuk mengelompokkan produk.

#### Fungsionalitas

| Fitur       | Deskripsi                                                          |
|-------------|--------------------------------------------------------------------|
| **Create**  | Tambah kategori baru (contoh: Kopi, Non-Kopi, Makanan, Snack)      |
| **Read**    | Lihat daftar semua kategori beserta jumlah produk di masing-masing |
| **Update**  | Edit nama kategori                                                  |
| **Delete**  | Hapus kategori (hanya jika tidak ada produk aktif yang terkait)     |

#### Acceptance Criteria

- [ ] Owner dapat membuat kategori baru dengan nama unik
- [ ] Daftar kategori menampilkan jumlah produk aktif yang terkait
- [ ] Owner dapat mengedit nama kategori dan perubahan terefleksi pada produk terkait
- [ ] Sistem mencegah penghapusan kategori yang masih memiliki produk aktif dan menampilkan pesan error yang jelas
- [ ] Nama kategori bersifat unik; sistem menolak duplikasi nama

---

### 3.3 Manajemen Stok

Owner dapat memonitor dan menyesuaikan stok produk, serta melihat riwayat pergerakan stok.

#### Fungsionalitas

| Fitur                     | Deskripsi                                                                          |
|---------------------------|------------------------------------------------------------------------------------|
| **Lihat Stok**            | Melihat jumlah stok saat ini per produk                                             |
| **Penyesuaian Stok**      | Tambah atau kurangi stok secara manual dengan alasan (restock, koreksi, kerusakan)  |
| **Riwayat Pergerakan Stok** | Log setiap perubahan stok: tipe (masuk/keluar), jumlah, alasan, waktu, user       |

#### Tipe Pergerakan Stok

| Tipe         | Kode        | Deskripsi                                           |
|--------------|-------------|-----------------------------------------------------|
| Stok Masuk   | `IN`        | Penambahan stok (restock dari supplier)              |
| Stok Keluar  | `OUT`       | Pengurangan stok (penjualan, kerusakan, expired)     |
| Koreksi      | `ADJUST`    | Penyesuaian manual hasil stock opname                |

#### Acceptance Criteria

- [ ] Owner dapat melihat stok saat ini untuk setiap produk
- [ ] Owner dapat menambah stok dengan mengisi jumlah dan alasan; stok bertambah sesuai input
- [ ] Owner dapat mengurangi stok dengan mengisi jumlah dan alasan; sistem menolak jika jumlah pengurangan melebihi stok saat ini
- [ ] Setiap perubahan stok tercatat di tabel riwayat pergerakan stok dengan informasi: produk, tipe, jumlah, alasan, timestamp, dan user yang melakukan
- [ ] Owner dapat melihat riwayat pergerakan stok per produk dengan filter berdasarkan tipe dan rentang tanggal

---

### 3.4 Manajemen Meja

Owner dapat mengelola daftar meja yang tersedia di coffee shop.

#### Fungsionalitas

| Fitur       | Deskripsi                                                            |
|-------------|----------------------------------------------------------------------|
| **Create**  | Tambah meja baru dengan nomor/nama meja dan kapasitas                |
| **Read**    | Lihat daftar meja beserta status saat ini (tersedia/terisi)          |
| **Update**  | Edit informasi meja (nomor, kapasitas)                               |
| **Delete**  | Hapus meja (hanya jika tidak ada transaksi aktif di meja tersebut)   |

#### Acceptance Criteria

- [ ] Owner dapat menambah meja baru dengan nomor meja yang unik
- [ ] Daftar meja menampilkan status real-time (tersedia/terisi)
- [ ] Owner dapat mengedit informasi meja
- [ ] Sistem mencegah penghapusan meja yang sedang digunakan dalam transaksi aktif
- [ ] Nomor meja bersifat unik; sistem menolak duplikasi

---

### 3.5 Manajemen User Cashier

Owner dapat mengelola akun-akun cashier yang mengoperasikan sistem.

#### Fungsionalitas

| Fitur               | Deskripsi                                                            |
|----------------------|----------------------------------------------------------------------|
| **Create**           | Buat akun cashier baru dengan nama, email, dan password              |
| **Read**             | Lihat daftar semua cashier beserta status aktif/nonaktif             |
| **Update**           | Edit informasi cashier (nama, email, reset password)                 |
| **Aktivasi/Deaktivasi** | Toggle status akun cashier; cashier nonaktif tidak bisa login     |

#### Acceptance Criteria

- [ ] Owner dapat membuat akun cashier baru dengan email unik dan password yang di-hash
- [ ] Daftar cashier menampilkan nama, email, status akun, dan tanggal pembuatan
- [ ] Owner dapat mengedit informasi cashier dan me-reset password
- [ ] Cashier dengan status nonaktif tidak dapat login ke sistem
- [ ] Owner tidak dapat menonaktifkan akunnya sendiri

---

### 3.6 Dashboard Laporan

Owner memiliki akses ke dashboard yang menyajikan ringkasan performa bisnis.

#### Fungsionalitas

| Fitur                        | Deskripsi                                                       |
|------------------------------|-----------------------------------------------------------------|
| **Revenue Harian**           | Total pendapatan hari ini                                        |
| **Revenue Mingguan**         | Total pendapatan 7 hari terakhir dengan grafik tren              |
| **Revenue Bulanan**          | Total pendapatan 30 hari terakhir dengan grafik tren             |
| **Produk Terlaris**          | Daftar top 10 produk berdasarkan jumlah terjual dalam periode    |
| **Transaksi per Cashier**    | Jumlah transaksi dan total revenue yang diproses per cashier     |

#### Acceptance Criteria

- [ ] Dashboard menampilkan revenue harian yang akurat berdasarkan transaksi berstatus `paid`
- [ ] Revenue mingguan dan bulanan menampilkan data agregat sesuai periode dengan visualisasi tren
- [ ] Produk terlaris menampilkan top 10 produk berdasarkan quantity terjual dalam periode yang dipilih
- [ ] Laporan transaksi per cashier menampilkan jumlah transaksi dan total revenue masing-masing cashier
- [ ] Semua data dashboard hanya menghitung transaksi yang berstatus `paid` (sukses)

---

### 3.7 Manajemen Promo dan Diskon

Owner dapat membuat dan mengelola promo/diskon yang bisa diaplikasikan ke transaksi.

#### Fungsionalitas

| Fitur              | Deskripsi                                                                |
|--------------------|--------------------------------------------------------------------------|
| **Create**         | Buat promo baru dengan nama, tipe diskon, nilai, periode aktif           |
| **Read**           | Lihat daftar promo beserta status (aktif/tidak aktif/expired)            |
| **Update**         | Edit informasi promo                                                      |
| **Delete**         | Hapus promo (soft delete)                                                 |
| **Tipe Diskon**    | Persentase (%) atau nominal tetap (Rp)                                   |

#### Struktur Data Promo

| Field                | Tipe       | Deskripsi                                       |
|----------------------|------------|-------------------------------------------------|
| `name`               | string     | Nama promo (contoh: "Diskon Weekend 10%")        |
| `discount_type`      | enum       | `percentage` atau `fixed`                        |
| `discount_value`     | decimal    | Nilai diskon (10 untuk 10% atau 5000 untuk Rp5000) |
| `min_purchase`       | decimal    | Minimum total belanja untuk memenuhi syarat promo |
| `max_discount`       | decimal    | Maksimum potongan (khusus tipe persentase)        |
| `start_date`         | datetime   | Tanggal mulai berlaku                             |
| `end_date`           | datetime   | Tanggal berakhir                                  |
| `is_active`          | boolean    | Status manual aktif/nonaktif                      |

#### Acceptance Criteria

- [ ] Owner dapat membuat promo baru dengan semua field wajib
- [ ] Promo bertipe persentase menghitung diskon berdasarkan persentase dari total belanja, dengan batas maksimum potongan (`max_discount`)
- [ ] Promo bertipe nominal memberikan potongan tetap sesuai `discount_value`
- [ ] Promo hanya bisa diapply jika: tanggal saat ini berada dalam rentang `start_date`–`end_date`, status `is_active` adalah `true`, dan total belanja memenuhi `min_purchase`
- [ ] Owner dapat mengedit dan melakukan soft delete pada promo
- [ ] Daftar promo menampilkan badge status: Aktif, Tidak Aktif, atau Expired

---

### 3.8 Export Laporan ke CSV

Owner dapat mengunduh data laporan dalam format CSV untuk kebutuhan analisis lebih lanjut.

#### Fungsionalitas

| Laporan                      | Deskripsi                                                  |
|------------------------------|------------------------------------------------------------|
| **Laporan Transaksi**        | Daftar transaksi dengan filter tanggal dan cashier          |
| **Laporan Revenue**          | Ringkasan revenue per hari dalam rentang tanggal tertentu   |
| **Laporan Produk Terlaris**  | Daftar produk berdasarkan jumlah terjual                    |
| **Laporan Stok**             | Riwayat pergerakan stok dalam rentang tanggal tertentu      |

#### Acceptance Criteria

- [ ] Owner dapat mengunduh laporan transaksi dalam format CSV dengan filter rentang tanggal dan cashier
- [ ] Owner dapat mengunduh laporan revenue harian dalam format CSV
- [ ] Owner dapat mengunduh laporan produk terlaris dalam format CSV
- [ ] Owner dapat mengunduh laporan pergerakan stok dalam format CSV
- [ ] File CSV memiliki header kolom yang jelas dan encoding UTF-8
- [ ] Nama file CSV mengandung tipe laporan dan rentang tanggal (contoh: `transaksi_2026-03-01_2026-03-30.csv`)

---

## 4. Fitur Cashier

### 4.1 Manajemen Shift

Cashier wajib membuka shift sebelum dapat memproses transaksi. Sistem shift memastikan akuntabilitas kas per cashier per hari.

#### Fungsionalitas

| Fitur             | Deskripsi                                                                      |
|-------------------|--------------------------------------------------------------------------------|
| **Buka Shift**    | Cashier membuka shift dengan memasukkan jumlah modal kas awal                   |
| **Tutup Shift**   | Cashier menutup shift dengan rekap: total transaksi, total revenue, kas akhir   |
| **Status Shift**  | Sistem menampilkan status shift saat ini (aktif/belum dibuka)                   |

#### Acceptance Criteria

- [ ] Cashier harus membuka shift sebelum bisa membuat transaksi; jika shift belum dibuka, tombol transaksi baru di-disable
- [ ] Saat buka shift, cashier wajib mengisi jumlah modal kas awal (angka positif)
- [ ] Satu cashier hanya bisa memiliki satu shift aktif dalam satu waktu
- [ ] Saat tutup shift, sistem menampilkan rekap: jumlah transaksi, total revenue, modal kas awal, dan estimasi kas akhir (modal + revenue)
- [ ] Setelah shift ditutup, cashier tidak bisa membuat transaksi baru sampai shift baru dibuka
- [ ] Data shift tercatat: cashier_id, waktu buka, waktu tutup, modal kas awal, total transaksi, total revenue

---

### 4.2 Buat Transaksi Baru

Cashier dapat membuat transaksi baru untuk melayani pelanggan.

#### Alur Transaksi

```
Pilih Meja → Tambah Produk ke Keranjang → Atur Quantity → (Opsional) Apply Promo → Checkout
```

#### Fungsionalitas

| Fitur                  | Deskripsi                                                                   |
|------------------------|-----------------------------------------------------------------------------|
| **Pilih Meja**         | Pilih meja yang statusnya tersedia; meja otomatis berubah jadi "terisi"     |
| **Pilih Produk**       | Browse daftar produk aktif, cari/filter berdasarkan kategori                |
| **Atur Quantity**      | Tambah/kurangi jumlah item; sistem validasi stok tersedia                   |
| **Keranjang**          | Review daftar item, subtotal per item, dan total keseluruhan                |
| **Hapus Item**         | Hapus item tertentu dari keranjang                                           |

#### Acceptance Criteria

- [ ] Cashier dapat memilih meja dari daftar meja yang berstatus tersedia
- [ ] Setelah meja dipilih dan transaksi dibuat, status meja otomatis berubah menjadi "terisi"
- [ ] Cashier dapat mencari dan memfilter produk berdasarkan kategori; hanya produk aktif yang tampil
- [ ] Cashier dapat menambah produk ke keranjang dan mengatur quantity
- [ ] Sistem menolak jika quantity melebihi stok yang tersedia dan menampilkan pesan error
- [ ] Keranjang menampilkan nama produk, harga satuan, quantity, subtotal per item, dan total keseluruhan
- [ ] Cashier dapat menghapus item dari keranjang

---

### 4.3 Apply Promo ke Transaksi

Cashier dapat menerapkan promo aktif yang tersedia ke dalam transaksi.

#### Fungsionalitas

| Fitur                | Deskripsi                                                           |
|----------------------|---------------------------------------------------------------------|
| **Lihat Promo Aktif**| Daftar promo yang sedang berlaku dan memenuhi syarat                |
| **Apply Promo**      | Terapkan satu promo ke transaksi; sistem hitung diskon otomatis     |
| **Hapus Promo**      | Batalkan promo yang sudah di-apply sebelum checkout                 |

#### Acceptance Criteria

- [ ] Cashier dapat melihat daftar promo yang aktif dan memenuhi syarat (`min_purchase`)
- [ ] Cashier dapat menerapkan satu promo per transaksi
- [ ] Sistem menghitung diskon secara otomatis sesuai tipe promo (persentase atau nominal) dan menampilkan total sebelum & sesudah diskon
- [ ] Untuk tipe persentase, diskon tidak melebihi `max_discount`
- [ ] Cashier dapat membatalkan promo yang sudah di-apply sebelum melakukan checkout
- [ ] Promo yang tidak memenuhi `min_purchase` tidak muncul di daftar atau ditandai sebagai "belum memenuhi syarat"

---

### 4.4 Checkout dengan Midtrans Snap

Proses pembayaran dilakukan melalui payment gateway Midtrans menggunakan Snap.

#### Alur Checkout

```
Review Keranjang → Klik Checkout → Generate Snap Token (Backend) → Tampilkan Snap Popup → Proses Pembayaran → Callback/Webhook → Update Status Transaksi
```

#### Fungsionalitas

| Fitur                   | Deskripsi                                                                |
|-------------------------|--------------------------------------------------------------------------|
| **Generate Snap Token** | Backend mengirim detail transaksi ke Midtrans API dan menerima snap token |
| **Snap Popup**          | Frontend menampilkan popup pembayaran Midtrans Snap                       |
| **Webhook Handler**     | Backend menerima notifikasi status pembayaran dari Midtrans               |
| **Update Status**       | Sistem update status transaksi berdasarkan notifikasi Midtrans            |

#### Status Transaksi

| Status       | Deskripsi                                            |
|--------------|------------------------------------------------------|
| `pending`    | Transaksi dibuat, menunggu pembayaran                 |
| `paid`       | Pembayaran berhasil dikonfirmasi oleh Midtrans        |
| `failed`     | Pembayaran gagal atau expired                         |
| `cancelled`  | Transaksi dibatalkan oleh cashier sebelum pembayaran  |

#### Acceptance Criteria

- [ ] Saat cashier klik Checkout, backend berhasil generate snap token dari Midtrans API
- [ ] Snap popup Midtrans tampil di frontend dan pelanggan dapat memilih metode pembayaran
- [ ] Setelah pembayaran berhasil, webhook Midtrans mengirim notifikasi ke backend dan status transaksi di-update menjadi `paid`
- [ ] Stok produk otomatis berkurang sesuai quantity setelah status menjadi `paid`
- [ ] Status meja otomatis kembali menjadi "tersedia" setelah transaksi berstatus `paid`, `failed`, atau `cancelled`
- [ ] Transaksi yang `pending` lebih dari batas waktu tertentu (sesuai konfigurasi Midtrans) otomatis menjadi `failed`
- [ ] Signature verification dari webhook Midtrans diimplementasikan untuk keamanan

---

### 4.5 Riwayat Transaksi Shift Hari Ini

Cashier dapat melihat daftar transaksi yang telah diproses selama shift aktif hari ini.

#### Fungsionalitas

| Fitur                   | Deskripsi                                                         |
|-------------------------|-------------------------------------------------------------------|
| **Daftar Transaksi**    | List semua transaksi dalam shift aktif dengan info ringkas         |
| **Detail Transaksi**    | Lihat detail lengkap per transaksi (item, qty, diskon, total, status) |

#### Acceptance Criteria

- [ ] Cashier dapat melihat daftar transaksi dalam shift aktif hari ini
- [ ] Daftar menampilkan: nomor transaksi, meja, waktu, total, dan status pembayaran
- [ ] Cashier dapat melihat detail transaksi termasuk daftar item, quantity, harga, diskon (jika ada), dan total akhir
- [ ] Cashier hanya dapat melihat transaksi milik shift-nya sendiri, bukan transaksi cashier lain

---

## 5. Business Rules

### Autentikasi & Otorisasi

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 1  | Sistem menggunakan JWT (JSON Web Token) untuk autentikasi                                           |
| 2  | Setiap request API harus menyertakan token JWT yang valid di header `Authorization`                  |
| 3  | Role `owner` dan `cashier` memiliki permission yang berbeda; endpoint di-protect berdasarkan role    |
| 4  | Password di-hash menggunakan bcrypt sebelum disimpan ke database                                     |
| 5  | Token JWT memiliki masa berlaku (expiry) yang dapat dikonfigurasi                                    |

### Transaksi

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 6  | Cashier wajib memiliki shift aktif sebelum membuat transaksi                                        |
| 7  | Setiap transaksi harus terkait dengan satu meja                                                     |
| 8  | Stok produk di-validasi saat menambah item ke keranjang; transaksi ditolak jika stok tidak cukup    |
| 9  | Pengurangan stok hanya dilakukan setelah webhook Midtrans mengkonfirmasi pembayaran berhasil (status `paid`), bukan saat transaksi dibuat |
| 10 | Satu transaksi hanya boleh menggunakan maksimal satu promo/diskon                                   |
| 11 | Harga produk yang disimpan di detail transaksi adalah harga saat transaksi dibuat (snapshot), bukan referensi ke harga produk saat ini |
| 12 | Order yang sudah masuk proses checkout (status `pending`, `paid`, `failed`) tidak bisa diubah item-nya (tambah, hapus, atau ubah quantity); perubahan hanya diperbolehkan selama transaksi masih dalam tahap keranjang (belum checkout) |

### Stok

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 13 | Stok tidak boleh bernilai negatif                                                                   |
| 14 | Setiap perubahan stok wajib tercatat di tabel riwayat pergerakan stok                               |
| 15 | Produk dengan stok 0 tetap tampil di daftar kasir tetapi tidak bisa ditambahkan ke keranjang         |

### Produk & Kategori

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 16 | Produk yang di-soft-delete atau berstatus nonaktif tidak tampil di daftar produk kasir               |
| 17 | Kategori tidak dapat dihapus jika masih memiliki produk aktif                                        |
| 18 | Nama produk dan nama kategori harus unik dalam scope masing-masing                                  |

### Meja

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 19 | Meja dengan transaksi aktif (status `pending`) tidak bisa dihapus atau dipilih untuk transaksi baru  |
| 20 | Status meja otomatis berubah: "tersedia" → "terisi" saat transaksi dibuat; "terisi" → "tersedia" saat transaksi selesai/gagal/dibatalkan |

### Shift

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 21 | Satu cashier hanya boleh memiliki satu shift aktif dalam satu waktu                                  |
| 22 | Shift tidak bisa ditutup jika masih ada transaksi berstatus `pending`                                |
| 23 | Data shift bersifat immutable setelah ditutup — tidak bisa diedit                                    |

### Promo & Diskon

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 24 | Promo hanya berlaku jika: `is_active = true`, tanggal saat ini dalam rentang `start_date`–`end_date`, dan total belanja ≥ `min_purchase` |
| 25 | Diskon tipe persentase dibatasi oleh `max_discount` (jika diatur)                                    |
| 26 | Diskon tipe nominal tidak boleh melebihi total belanja (total setelah diskon minimal Rp 0)           |

### Pembayaran

| No | Rule                                                                                               |
|----|-----------------------------------------------------------------------------------------------------|
| 27 | Semua pembayaran diproses melalui Midtrans Snap; sistem tidak menangani pembayaran tunai di versi ini |
| 28 | Webhook dari Midtrans harus diverifikasi menggunakan signature key sebelum diproses                  |
| 29 | Status transaksi hanya boleh di-update melalui webhook Midtrans, bukan manual                        |

---

## 6. Out of Scope

Fitur-fitur berikut **tidak akan dibangun** pada versi 1.0 ini:

| Fitur                             | Alasan                                                         |
|-----------------------------------|----------------------------------------------------------------|
| **Pembayaran Tunai (Cash)**       | Fokus versi ini pada pembayaran digital via Midtrans            |
| **Multi-outlet / Multi-cabang**   | Sistem dirancang untuk satu outlet saja                         |
| **Manajemen Supplier**            | Pengelolaan supplier di luar scope POS                          |
| **Kitchen Display System (KDS)**  | Integrasi dapur bisa ditambahkan di versi mendatang             |
| **Loyalty Program / Membership**  | Fitur loyalitas pelanggan belum diprioritaskan                  |
| **Inventory Forecasting**         | Prediksi stok berbasis AI di luar scope versi ini               |
| **Mobile App (iOS/Android)**      | Fokus pada web-based application terlebih dahulu                |
| **Cetak Struk / Printer Integration** | Integrasi printer thermal belum disupport                   |
| **Manajemen Bahan Baku / Resep** | Hanya mengelola stok produk jadi, bukan bahan baku              |
| **Multi-bahasa (i18n)**           | Sistem menggunakan Bahasa Indonesia saja                        |
| **Notifikasi Real-time (WebSocket)** | Notifikasi push belum diimplementasikan                      |
| **Audit Log Komprehensif**        | Hanya riwayat stok yang di-log; audit log umum di versi mendatang |
| **Split Bill / Pisah Tagihan**    | Satu transaksi satu pembayaran utuh                             |
| **Reservasi Meja**                | Fitur reservasi online tidak termasuk                           |

---

> [!NOTE]
> Dokumen ini bersifat **living document** dan akan diperbarui seiring perkembangan proyek. Setiap perubahan signifikan harus melalui review dan approval stakeholder terkait.
