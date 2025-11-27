# Implementasi Login & Autentikasi (FR-001)

## Deskripsi
Implementasi login menggunakan username/email dan password dengan JWT token yang berisi role dan permissions.

## Struktur File

### Models (`app/models/`)
- **user.go**: Struct User dengan field Username, Email, PasswordHash, IsActive
- **role.go**: Struct Role
- **permission.go**: Struct Permission
- **role_permission.go**: Struct RolePermission
- **dto.go**: Data Transfer Objects (LoginRequest, LoginResponse, UserProfile, RoleProfile, JWTClaims)

### Repository (`app/repository/`)

#### user_repository.go
Query database untuk User:
- `FindByUsername()`: Cari user berdasarkan username
- `FindByEmail()`: Cari user berdasarkan email
- `FindByID()`: Cari user berdasarkan ID
- `GetUserWithRoleAndPermissions()`: Ambil user beserta role dan permissions

#### role_permission_repository.go
Query database untuk Role dan Permission:
- `FindByName()`: Cari role/permission berdasarkan nama
- `AssignPermissionToRole()`: Assign permission ke role
- `GetPermissionsByRole()`: Ambil semua permissions untuk role

### Service (`app/service/`)

#### auth_service.go
Business logic untuk autentikasi:
- `Login()`: Handle login flow
  1. Validasi input (username/email dan password)
  2. Cari user di database (FindByUsername atau FindByEmail)
  3. Cek status aktif user (IsActive == true)
  4. Verifikasi password (VerifyPassword)
  5. Ambil user dengan role dan permissions (GetUserWithRoleAndPermissions)
  6. Generate JWT token dengan role dan permissions
  7. Return LoginResponse dengan token dan user profile

### Utils (`utils/`)

#### jwt_utils.go
Utility functions:
- `HashPassword()`: Hash password menggunakan SHA256
- `VerifyPassword()`: Verifikasi password dengan hash
- `GenerateJWT()`: Generate JWT token dengan claims (user_id, username, email, role, permissions, exp, iat)

## Flow Login (FR-001)

```
┌─────────────────────────────────────┐
│ 1. User Input (username/email+pass) │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 2. AuthService.Login() Validasi     │
└────────────┬────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│ 3. UserRepo.FindByUsername/Email()   │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│ 4. Check IsActive == true            │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│ 5. VerifyPassword() - SHA256 compare │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────────────────┐
│ 6. GetUserWithRoleAndPermissions()                   │
│    - Ambil user dengan role                          │
│    - Query permissions via role_permission table     │
└────────────┬─────────────────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│ 7. GenerateJWT() dengan claims       │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│ 8. Return LoginResponse              │
│    - Token                           │
│    - UserProfile (tanpa password)    │
└──────────────────────────────────────┘
```

## Model User

```go
type User struct {
    ID           string    `gorm:"primaryKey"`
    Username     string    `gorm:"uniqueIndex"`
    Email        string    `gorm:"uniqueIndex"`
    PasswordHash string    // SHA256 hashed password
    FullName     string
    RoleID       string
    IsActive     bool      `gorm:"default:true"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

## JWT Claims

```json
{
  "user_id": "user-uuid-string",
  "username": "admin",
  "email": "admin@example.com",
  "role": "Admin",
  "permissions": ["create_user", "edit_user", "delete_user"],
  "exp": 1234567890,
  "iat": 1234567800
}
```

## Cara Menggunakan

### 1. Inisialisasi Services (di main.go)

```go
userRepo := repository.NewUserRepository()
authService := service.NewAuthService(userRepo)
```

### 2. Login User

```go
loginReq := &models.LoginRequest{
    Username: "admin",
    Password: "password123",
}

response, err := authService.Login(loginReq)
if err != nil {
    log.Printf("Login failed: %v", err)
    return
}

log.Printf("Token: %s", response.Token)
log.Printf("User: %+v", response.Profile)
```

### 3. Create User Baru

```go
hashedPassword := utils.HashPassword("password123")

user := &models.User{
    Username:     "admin",
    Email:        "admin@example.com",
    FullName:     "Admin User",
    PasswordHash: hashedPassword,
    IsActive:     true,
    RoleID:       "role-id-from-db",
}

database.DB.Create(user)
```

## Environment Variables

Tambahkan ke `.env`:
```
JWT_SECRET=your-secret-key-change-in-production
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=database_name
```

Jika JWT_SECRET tidak ada, akan menggunakan default key (jangan di production!)

## Request/Response Example

### Login Request
```json
{
  "username": "admin",
  "password": "password123"
}
```

atau

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

### Login Response (Success)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "profile": {
    "id": "user-uuid",
    "full_name": "Admin User",
    "email": "admin@example.com",
    "username": "admin",
    "is_active": true,
    "role": {
      "id": "role-uuid",
      "name": "Admin"
    },
    "permissions": [
      {"id": "1", "name": "create_user", "resource": "User", "action": "create"},
      {"id": "2", "name": "edit_user", "resource": "User", "action": "edit"},
      {"id": "3", "name": "delete_user", "resource": "User", "action": "delete"}
    ]
  }
}
```

### Login Response (Error)
```json
{
  "error": "invalid credentials"
}
```

## Catatan Penting

✅ **Tidak ada HTTP Handler/Controller** - hanya menggunakan Service & Repository
✅ **Query logic di Repository** - semua database query
✅ **Business logic di Service** - validation, authentication flow
✅ **Model hanya struct** - hanya definisi struktur data
✅ **Password di-hash** - menggunakan SHA256 sebelum disimpan
✅ **JWT berisi role dan permissions** - untuk authorization
✅ **Validasi status aktif user** - hanya active user yang bisa login
✅ **Support login dengan username atau email** - fleksibel
✅ **Permissions di-fetch dari database** - via JOIN query

## Testing

Bisa ditest menggunakan:
- Postman: Import JWT token di Authorization header
- Unit tests untuk AuthService
- Integration tests dengan database

Contoh di `example_login.go`

