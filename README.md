# Sistem Manajemen Prestasi Mahasiswa

API Backend untuk sistem pencatatan dan manajemen prestasi mahasiswa dengan fitur verifikasi dari dosen pembimbing akademik.

## Tentang Project

Project ini adalah sistem backend untuk mengelola data prestasi mahasiswa di perguruan tinggi. Sistemnya dibuat dengan fokus pada keamanan dan kontrol akses berbasis peran (Role-Based Access Control). 

Fitur utamanya:
- Mahasiswa bisa submit prestasi mereka (lomba, kompetisi, karya ilmiah, dll)
- Dosen Wali bisa verifikasi atau reject prestasi dari mahasiswa bimbingannya
- Admin bisa kelola semua data user, mahasiswa, dan dosen
- Setiap role punya akses yang berbeda-beda sesuai kebutuhannya

## Tech Stack

**Backend Framework:**
- Go 1.21+
- Fiber v2 (web framework yang cepat, mirip Express.js)

**Database:**
- PostgreSQL (data relational: users, students, lecturers, dll)
- MongoDB (data prestasi yang lebih kompleks dengan nested objects)

**Authentication & Security:**
- JWT (JSON Web Token) untuk autentikasi
- Access token: 1 jam
- Refresh token: 7 hari
- RBAC (Role-Based Access Control)

**Documentation:**
- Swagger/OpenAPI (auto-generated dari code comments)

**Library Penting:**
- GORM (ORM untuk PostgreSQL)
- MongoDB Go Driver
- golang-jwt/jwt
- go-playground/validator

## Instalasi

### Prerequisites

Pastikan sudah install:
- Go versi 1.21 atau lebih baru
- PostgreSQL
- MongoDB
- Git

### Setup Database

**PostgreSQL:**
```sql
CREATE DATABASE db_uas;
```

**MongoDB:**
Buat database baru dengan nama `db_uas`

### Install Dependencies

```bash
go mod download
```

### Environment Variables

Buat file `.env` di root folder:

```env
# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=db_uas

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=db_uas

# JWT Secret
JWT_SECRET=your-super-secret-key-change-this-in-production

# Server
PORT=8080
```

### Jalankan Aplikasi

```bash
# Build aplikasi
go build

# Atau langsung run
go run main.go
```

Server akan jalan di `http://localhost:8080`

### Generate Swagger Documentation

Kalau ada perubahan di API:

```bash
swag init --parseInternal --parseDependency
```

## Struktur Project

```
.
├── app/
│   ├── models/          # Definisi struct/model database
│   ├── repository/      # Database access layer
│   └── service/         # Business logic & handlers
├── config/              # Konfigurasi aplikasi
├── database/            # Setup database & migrations
├── docs/                # Swagger documentation (auto-generated)
├── middleware/          # Auth, RBAC, Ownership middleware
├── routes/              # Route definitions
├── utils/               # Helper functions (JWT, response formatter)
├── uploads/             # File uploads (achievement attachments)
├── main.go              # Entry point
└── go.mod               # Dependencies
```

## Role & Permissions

Ada 4 role dalam sistem:

### 1. Admin
Akses penuh ke semua endpoint. Bisa:
- Kelola semua user (create, read, update, delete)
- Lihat dan kelola semua data mahasiswa
- Lihat dan kelola semua data dosen
- Lihat semua prestasi
- Assign dosen wali ke mahasiswa

### 2. Mahasiswa
Akses terbatas ke data sendiri. Bisa:
- Submit prestasi baru (draft)
- Edit prestasi yang masih draft
- Submit prestasi untuk verifikasi
- Lihat prestasi sendiri
- Lihat profil sendiri

### 3. Dosen Wali
Akses berikut ini. Bisa:
- Lihat profil dosen sendiri
- Lihat mahasiswa bimbingan sendiri
- Lihat prestasi mahasiswa bimbingan
- Verifikasi atau reject prestasi mahasiswa bimbingan
- Lihat statistik prestasi mahasiswa bimbingan

## API Documentation

### Cara Akses Swagger UI

Buka browser dan akses:
```
http://localhost:8080/swagger/index.html
```

## API Endpoints

### Authentication

```
POST   /api/v1/auth/login           # Login dan dapat token
POST   /api/v1/auth/logout          # Logout (client-side)
POST   /api/v1/auth/refresh         # Refresh access token
GET    /api/v1/auth/profile         # Lihat profil user login
```

### User Management (Admin only)

```
POST   /api/v1/users                # Buat user baru
GET    /api/v1/users                # List semua user
GET    /api/v1/users/:id            # Detail user
PUT    /api/v1/users/:id            # Update user
DELETE /api/v1/users/:id            # Hapus user (hard delete)
PUT    /api/v1/users/:id/role       # Ganti role user
```

### Students

```
GET    /api/v1/students                     # List mahasiswa
POST   /api/v1/students                     # Buat profil mahasiswa (Admin)
GET    /api/v1/students/:id                 # Detail mahasiswa
PUT    /api/v1/students/:id                 # Update profil (Admin)
GET    /api/v1/students/:id/achievements    # Prestasi mahasiswa
PUT    /api/v1/students/:id/advisor         # Set dosen wali (Admin)
```

### Lecturers

```
GET    /api/v1/lecturers                # List dosen
POST   /api/v1/lecturers                # Buat profil dosen (Admin)
PUT    /api/v1/lecturers/:id            # Update profil (Admin)
GET    /api/v1/lecturers/:id/advisees   # Mahasiswa bimbingan
```

### Achievements

```
GET    /api/v1/achievements              # List prestasi
POST   /api/v1/achievements              # Buat prestasi baru
GET    /api/v1/achievements/:id          # Detail prestasi
PUT    /api/v1/achievements/:id          # Update prestasi
DELETE /api/v1/achievements/:id          # Hapus prestasi
POST   /api/v1/achievements/:id/submit   # Submit untuk verifikasi
POST   /api/v1/achievements/:id/verify   # Verifikasi (Dosen Wali)
POST   /api/v1/achievements/:id/reject   # Reject (Dosen Wali)
POST   /api/v1/achievements/:id/upload   # Upload lampiran
GET    /api/v1/achievements/:id/history  # History perubahan
```

### Reports & Statistics

```
GET    /api/v1/reports/statistics           # Statistik prestasi
GET    /api/v1/reports/student/:student_id  # Report mahasiswa
```

## Workflow Prestasi

### 1. Mahasiswa Submit Prestasi

**Step 1: Buat draft**
```
POST /api/v1/achievements
{
  "title": "Juara 1 Lomba Programming",
  "achievement_type": "competition",
  "description": "...",
  "details": {...}
}
```

Status: `draft`

**Step 2: Submit untuk verifikasi**
```
POST /api/v1/achievements/{id}/submit
```

Status berubah jadi: `submitted`

### 2. Dosen Wali Verifikasi

**Approve:**
```
POST /api/v1/achievements/{id}/verify
```

Status berubah jadi: `verified`

**Reject:**
```
POST /api/v1/achievements/{id}/reject
{
  "rejection_note": "Bukti kurang lengkap"
}
```

Status berubah jadi: `rejected`

## Keamanan & Access Control

### Authentication

Semua endpoint (kecuali login) butuh JWT token di header:
```
Authorization: Bearer <access_token>
```

Token otomatis di-validate oleh middleware sebelum request sampai ke handler.

### Authorization (RBAC)

Sistem pakai permission-based authorization. Setiap endpoint punya requirement permission tertentu:

**Contoh:**
- `user:manage` - Kelola user (Admin only)
- `student:read` - Baca data mahasiswa
- `lecturer:read` - Baca data dosen
- `achievement:read` - Baca prestasi
- `achievement:create` - Buat prestasi
- `achievement:verify` - Verifikasi prestasi (Dosen Wali)

Permission sudah dibuat sebelumnya di databse potsgreysql.

### Ownership Validation

Selain role dan permission, ada validasi kepemilikan data:

**Mahasiswa:**
- Cuma bisa akses prestasi sendiri
- Gak bisa lihat prestasi mahasiswa lain

**Dosen Wali:**
- Cuma bisa lihat mahasiswa bimbingannya
- Cuma bisa verifikasi prestasi dari mahasiswa bimbingannya 
- Gak bisa akses data mahasiswa/dosen lain

**Admin:**
- Full access tanpa batasan

## Testing

Run unit tests:
```bash
go test ./... -v
```

Run specific test:
```bash
go test ./app/service -v -run TestAchievementService
```

## Troubleshooting

### Error: "user not authenticated"
Pastikan token ada di header dan formatnya benar:
```
Authorization: Bearer <token>
```

### Error: "insufficient permissions"
User tidak punya permission untuk akses endpoint tersebut. Cek role dan permission user.

### Error: "you can only access your own advisees"
Dosen Wali coba akses data mahasiswa yang bukan bimbingannya.

### Database connection error
Cek file `.env` dan pastikan credentials database sudah benar.

## Default Users

Setelah pertama kali run, ada user default yang bisa dipakai:

**Admin:**
```
username: admin
password: admin123
```

Untuk production, segera ganti password default dan buat user baru.

## Development Notes

### Tambah Endpoint Baru

1. Tambah method di service interface
2. Implementasi method di service impl
3. Tambah route di routes/
4. Tambah Swagger comments
5. Generate ulang swagger: `swag init --parseInternal --parseDependency`

### Tambah Permission Baru

1. Tambah di `database/seed.go`
2. Assign ke role yang sesuai
3. Run ulang aplikasi (akan auto seed)

### Database Schema Changes

1. Ubah model di `app/models/`
2. Update repository kalau perlu
3. Drop database atau manual migration
4. Run ulang aplikasi (auto migrate)

## Production Deployment

Sebelum deploy ke production:

1. Ganti `JWT_SECRET` dengan random string yang kuat
2. Ganti password default admin
3. Setup proper database backup
4. Pakai HTTPS untuk semua koneksi
5. Setup rate limiting
6. Enable CORS dengan whitelist domain
7. Review dan audit permissions

## License

Project ini dibuat untuk keperluan akademik yaitu UAS Backend Praktikum .
