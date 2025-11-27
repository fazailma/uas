# Routes & Middleware Documentation

## 🔐 Authentication Routes (v1)

### 1. Login Endpoint
**POST** `/api/v1/auth/login`

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

atau menggunakan email:
```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

**Response (Success - 200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": "user-uuid",
  "message": "Login successful"
}
```

**Response (Error - 401):**
```json
{
  "error": "invalid credentials"
}
```

---

### 2. Refresh Token Endpoint
**POST** `/api/v1/auth/refresh`

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (Success - 200):**
```json
{
  "message": "token refreshed",
  "token": "new-token-here"
}
```

---

### 3. Logout Endpoint
**POST** `/api/v1/auth/logout`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (Success - 200):**
```json
{
  "message": "logout successful"
}
```

---

### 4. Get Profile Endpoint (Protected)
**GET** `/api/v1/auth/profile`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (Success - 200):**
```json
{
  "user_id": "user-uuid",
  "username": "admin",
  "email": "admin@example.com",
  "role": "Admin",
  "permissions": ["create_user", "edit_user", "delete_user"]
}
```

**Response (Error - 401):**
```json
{
  "error": "invalid or expired token"
}
```

---

## 🛡️ Middleware

### AuthMiddleware
Located: `middleware/auth_middleware.go`

**Fungsi:**
- Memvalidasi JWT token dari Authorization header
- Mengecek format "Bearer <token>"
- Memverifikasi signature token
- Menyimpan user claims di context untuk digunakan handler

**Penggunaan:**
```go
protected := app.Group("/api", middleware.AuthMiddleware)
protected.Get("/profile", GetProfileHandler)
```

**Error Responses:**
- 401: Missing authorization header
- 401: Invalid authorization header format
- 401: Invalid or expired token

---

## 📂 File Structure

```
routes/
├── auth_routes.go      // Handler untuk login, logout, profile
└── routes.go          // Setup semua routes

middleware/
└── auth_middleware.go // JWT validation middleware
```

---

## 🔄 Route Flow

```
┌─────────────────────────────────────┐
│ Request ke /api/v1/auth/login       │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ LoginHandler                        │
│ - Parse LoginCredential             │
│ - Call AuthService.Login()          │
│ - Return LoginResponse              │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Response dengan JWT Token           │
└─────────────────────────────────────┘


┌─────────────────────────────────────┐
│ Request ke /api/v1/auth/profile     │
│ Header: Authorization: Bearer <token>
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ AuthMiddleware                      │
│ - Extract token dari header         │
│ - Validate JWT signature            │
│ - Store claims di context           │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ GetProfileHandler                   │
│ - Get user data dari context        │
│ - Return user profile               │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ Response dengan user profile        │
└─────────────────────────────────────┘
```

---

## 🚀 Testing dengan Postman

### 1. Login
- Method: `POST`
- URL: `http://localhost:8080/api/v1/auth/login`
- Body (JSON):
  ```json
  {
    "username": "admin",
    "password": "password123"
  }
  ```
- Copy token dari response

### 2. Refresh Token
- Method: `POST`
- URL: `http://localhost:8080/api/v1/auth/refresh`
- Body (JSON):
  ```json
  {
    "token": "<token-dari-login>"
  }
  ```

### 3. Get Profile
- Method: `GET`
- URL: `http://localhost:8080/api/v1/auth/profile`
- Headers:
  - Key: `Authorization`
  - Value: `Bearer <token-dari-login>`

### 4. Logout
- Method: `POST`
- URL: `http://localhost:8080/api/v1/auth/logout`
- Headers:
  - Key: `Authorization`
  - Value: `Bearer <token-dari-login>`

---

## ⚙️ Konfigurasi

Di `main.go`:
```go
routes.SetupRoutes(app)
```

Routes yang didaftarkan:
1. **Public Routes (No Auth Required)**
   - POST `/api/auth/login`
   - POST `/api/auth/logout`

2. **Protected Routes (Require JWT)**
   - GET `/api/profile`
   - (Bisa ditambah endpoint lainnya)

---

## 🔑 JWT Claims

Token yang di-generate berisi claims:
```json
{
  "user_id": "string",
  "username": "string",
  "email": "string",
  "role": "string",
  "permissions": ["string"],
  "exp": unix_timestamp,
  "iat": unix_timestamp
}
```

Token berlaku selama **24 jam**.

---

## 📝 Menambah Route Baru

Untuk menambah route baru yang protected:

```go
// Di routes/routes.go
protected := app.Group("/api", middleware.AuthMiddleware)
protected.Get("/profile", GetProfileHandler)
protected.Post("/users", CreateUserHandler)  // Route baru
protected.Get("/users/:id", GetUserHandler)  // Route baru
```

Untuk menambah route baru yang public:

```go
// Di routes/routes.go
auth := app.Group("/api/auth")
auth.Post("/login", LoginHandler)
auth.Post("/register", RegisterHandler)  // Route baru
```
