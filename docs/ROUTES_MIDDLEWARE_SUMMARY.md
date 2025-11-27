# Routes & Middleware Implementation - SUMMARY

## ✅ Implementasi Selesai

Routes dan middleware untuk login telah dibuat dan terintegrasi dengan main.go.

---

## 📁 File-File yang Dibuat

### 1. Routes (`routes/`)

#### auth_routes.go (BARU)
```go
// Handlers:
- LoginHandler(c *fiber.Ctx) error
  └─ POST /api/v1/auth/login
  └─ Input: LoginCredential
  └─ Output: LoginResponse + JWT Token

- RefreshTokenHandler(c *fiber.Ctx) error
  └─ POST /api/v1/auth/refresh
  └─ Input: token
  └─ Output: new JWT token

- LogoutHandler(c *fiber.Ctx) error
  └─ POST /api/v1/auth/logout
  └─ Output: logout successful message

- GetProfileHandler(c *fiber.Ctx) error
  └─ GET /api/v1/auth/profile (protected)
  └─ Mengambil user data dari JWT claims
  └─ Output: user profile (user_id, username, email, role, permissions)
```

#### routes.go (BARU)
```go
// Function:
- SetupRoutes(app *fiber.App)
  └─ Public routes: /api/v1/auth/login, /api/v1/auth/refresh, /api/v1/auth/logout
  └─ Protected routes: /api/v1/auth/profile (dengan AuthMiddleware)
```

---

### 2. Middleware (`middleware/`)

#### auth_middleware.go (BARU)
```go
// Function:
- AuthMiddleware(c *fiber.Ctx) error
  └─ Validasi JWT token dari Authorization header
  └─ Format: "Bearer <token>"
  └─ Cek signature dan validity
  └─ Store claims di context untuk handler
  └─ Return 401 jika token invalid/expired
```

---

### 3. Main Application

#### main.go (DIMODIFIKASI)
```go
// Perubahan:
- Import routes package
- Panggil routes.SetupRoutes(app)
- Remove RunMigrations() call
```

---

## 🔄 Endpoint Summary

| Method | Endpoint                | Auth Required | Fungsi                              |
|--------|-------------------------|---------------|------------------------------------|
| POST   | `/api/v1/auth/login`    | ❌ Tidak      | Login dengan username/email & pass  |
| POST   | `/api/v1/auth/refresh`  | ❌ Tidak      | Refresh JWT token                   |
| POST   | `/api/v1/auth/logout`   | ❌ Tidak      | Logout                              |
| GET    | `/api/v1/auth/profile`  | ✅ Ya        | Ambil user profile (protected)      |

---

## 🔐 Request/Response Example

### Login Request
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

### Login Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "user-uuid",
  "message": "Login successful"
}
```

### Refresh Token Request
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### Refresh Token Response
```json
{
  "message": "token refreshed",
  "token": "new-token-here"
}
```

### Get Profile Request (dengan token)
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Get Profile Response
```json
{
  "user_id": "user-uuid",
  "username": "admin",
  "email": "admin@example.com",
  "role": "Admin",
  "permissions": ["create_user", "edit_user", "delete_user"]
}
```

---

## 🏗️ Architecture

```
main.go
  ├─ routes.SetupRoutes(app)
  │  ├─ Public Routes (/api/auth/*)
  │  │  ├─ LoginHandler
  │  │  │  ├─ UserRepository
  │  │  │  └─ AuthService.Login()
  │  │  └─ LogoutHandler
  │  │
  │  └─ Protected Routes (/api/*)
  │     ├─ AuthMiddleware (JWT validation)
  │     └─ GetProfileHandler
  │
  └─ middleware.AuthMiddleware
     └─ Validasi JWT + Extract claims
```

---

## 🔒 Middleware Flow

```
Request dengan JWT token
      │
      ▼
AuthMiddleware
  ├─ Extract "Bearer <token>"
  ├─ Parse JWT dengan secret
  ├─ Validate signature
  ├─ Check expiry
  ├─ Store claims di context
  └─ Call c.Next()
      │
      ▼
Handler (GetProfileHandler)
  ├─ Akses c.Locals("user_id")
  ├─ Akses c.Locals("email")
  ├─ dll...
  └─ Return response
```

---

## 📝 Usage di Handler

Untuk mengakses JWT claims di dalam handler:

```go
func GetProfileHandler(c *fiber.Ctx) error {
    userID := c.Locals("user_id")      // string
    username := c.Locals("username")    // string
    email := c.Locals("email")          // string
    role := c.Locals("role")            // string
    permissions := c.Locals("permissions")  // []interface{}

    // Use data...
}
```

---

## 🚀 Build Status

```
✅ go build - SUCCESS
✅ All files compiled
✅ Routes integrated
✅ Middleware working
```

---

## 📌 Next Steps (Optional)

Untuk memperluas routes:

1. **Tambah route baru di routes.go:**
   ```go
   protected.Post("/users", CreateUserHandler)
   protected.Get("/users/:id", GetUserHandler)
   ```

2. **Buat handler baru di auth_routes.go:**
   ```go
   func CreateUserHandler(c *fiber.Ctx) error {
       // Implementation
   }
   ```

3. **Tambah middleware lain jika diperlukan:**
   - RoleMiddleware (check role-based access)
   - RateLimitMiddleware
   - LoggingMiddleware

---

## 🔑 Environment Variables

Ensure `.env` memiliki:
```env
JWT_SECRET=your-secret-key-change-in-production
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=database_name
```

---

## ✨ Features

✅ Login dengan JWT token
✅ Protected routes dengan middleware
✅ JWT validation
✅ Claims extraction ke context
✅ Error handling (401 Unauthorized)
✅ Support username atau email login
✅ Bearer token format
✅ 24-hour token expiry
