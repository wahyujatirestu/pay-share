# Pay-Share

**Pay-Share** adalah aplikasi backend RESTful yang dibangun dengan bahasa **Go**. Aplikasi ini dirancang untuk menangani manajemen user, produk, dan transaksi jual beli, serta dilengkapi dengan integrasi pembayaran melalui **Midtrans** untuk berbagai metode pembayaran.

---

## 🚀 Fitur Utama

-   Manajemen **User** (Register, Login, Update, Delete) dengan peran `customer` dan `employee`
-   Manajemen **Produk** (CRUD)
-   **Transaksi Jual Beli** lengkap dengan detail produk
-   Integrasi pembayaran dengan **Midtrans**
-   Mendukung berbagai metode pembayaran: bank transfer, QRIS, GoPay, ShopeePay, kartu kredit, dll
-   **JWT Authentication** & **Refresh Token** untuk manajemen sesi pengguna
-   Hashing password aman dengan **Argon2**
-   Validasi input dengan **DTO**
-   Proteksi route dengan **middleware role-based**

---

## Teknologi & Library

| Komponen          | Teknologi                            |
| ----------------- | ------------------------------------ |
| Bahasa            | Go                                   |
| Framework Web     | Gin                                  |
| Database          | PostgreSQL                           |
| PostgreSQL Driver | `github.com/lib/pq`                  |
| UUID              | `github.com/google/uuid`             |
| Hashing           | Argon2 (`x/crypto/argon2`)           |
| Auth              | JWT (`github.com/golang-jwt/jwt/v5`) |
| HTTP Client       | Resty                                |
| Pembayaran        | Midtrans                             |

---

## Struktur Project

```
pay-share/
├── main.go
├── server.go
├── go.mod
├── go.sum
├── .env
├── .gitignore
├── api.rest                # Kumpulan endpoint API siap test
├── config/                 # Konfigurasi DB & ENV
├── controller/             # Handler HTTP
├── dto/                    # Data Transfer Object
├── model/                  # Representasi tabel DB
├── middleware/             # Middleware (JWT, role, dll)
├── payment/                # Integrasi Midtrans
├── repository/             # Akses data (CRUD)
├── routes/                 # Definisi route
├── service/                # Business logic
├── utils/                  # Helper, security, token
├── sql/                    # SQL schema
└── README.md
```

---

## Setup & Instalasi

1. **Clone repository**

    ```bash
    git clone https://github.com/wahyujatirestu/pay-share.git
    cd pay-share
    ```

2. **Install dependencies**

    ```bash
    go mod tidy
    ```

3. **Setup database**

    - Jalankan PostgreSQL
    - Jalankan file `sql/ddl.sql` untuk setup schema

4. **Buat file `.env`**

    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_user
    DB_PASSWORD=your_password
    DB_NAME=payshare
    DB_DRIVER=postgres
    JWT_SECRET=your_jwt_secret

    MIDTRANS_SERVER_KEY=your_midtrans_server_key
    ```

5. **Jalankan aplikasi**

    ```bash
    go run .
    ```

6. **Akses API**

    ```
    http://localhost:8080/api/v1
    ```

---

## Tentang Refresh Token

**Refresh Token** digunakan untuk memperpanjang masa berlaku akses token (JWT) tanpa harus login ulang. Ini memberikan keamanan lebih dan pengalaman pengguna yang mulus, karena:

-   Access token bisa kedaluwarsa dalam waktu singkat
-   Refresh token memungkinkan regenerasi token baru selama sesi masih valid

---

## Dokumentasi API

-   File `api.rest` berisi kumpulan endpoint siap diuji
-   Bisa digunakan dengan:

    -   VSCode REST Client Extension(Instal ekstensi REST Client terlebih dahulu jika belum terinstal)
    -   Postman

---

## Lisensi

**MIT License**
© 2025 Restu Adi Wahyujati

---

> Dibuat dengan semangat & secangkir susu jahe ☕️ oleh [Restu Adi Wahyujati](https://github.com/wahyujatirestu)
