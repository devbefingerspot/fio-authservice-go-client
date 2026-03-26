# fio-auth-service-go-client

Go client library untuk Fingerspot Auth Service. Menyediakan fungsi login, verifikasi JWT, OTP, S2S token, manajemen user-company, dan gRPC client untuk query internal.

---

## Instalasi

```bash
go get github.com/devbefingerspot/fio-authservice-go-client
```

> Sesuaikan module path dengan konfigurasi `go.mod` di project kamu.

---

## Inisialisasi Client

```go
import authclient "github.com/devbefingerspot/fio-authservice-go-client"

client := authclient.NewFioAuthClient(
    "http://localhost:8080",       // base URL auth service (HTTP)
    "auth-grpc.example.com:50051", // base URL gRPC server (kosongkan untuk pakai host yang sama)
    "my-api-key",                  // API key untuk gRPC (kosongkan jika tidak dipakai)
    30*time.Second,                // HTTP timeout
    // optional: cache TTL untuk JWKS (default 5 menit)
    // 10*time.Minute,
)
defer client.Close() // tutup koneksi gRPC saat selesai
```

### Opsi Tambahan

```go
// Nonaktifkan TLS pada koneksi gRPC (hanya untuk development/local)
client.WithGRPCInsecure()
```

---

## Contoh Penggunaan

### Health Check

```go
resp, err := client.HealthCheck()
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Message)
```

---

### Login Web

```go
resp, err := client.WebLogin("user@example.com", "password123", authclient.PlatformNewWeb)
if err != nil {
    if err.Error() == "invalid_credentials" {
        log.Println("Email atau password salah")
    }
    log.Fatal(err)
}

if resp.IsRedirect() {
    // Platform mismatch — user diarahkan ke platform lain
    fmt.Println("Redirect ke:", resp.RedirectPlatform)
    fmt.Println("OTC Token:", *resp.OTCToken)
} else {
    fmt.Println("Access Token:", resp.AccessToken)
    fmt.Println("Refresh Token:", resp.RefreshToken)
}
```

---

### Login Mobile

```go
// Login dengan email
resp, err := client.MobileLogin("password123", "user@example.com", "", "")

// Login dengan nomor telepon
resp, err := client.MobileLogin("password123", "", "81234567890", "62")

if err != nil {
    log.Fatal(err)
}
fmt.Println("Identity Access Token:", resp.IdentityAccessToken)

// Tukar identity token dengan company-scoped token
companyResp, err := client.MobileIssueCompanyToken(
    resp.IdentityAccessToken,
    "company-uuid",
    string(authclient.RoleEmployee),
    "device-uuid",
    "fcm-token",
    "MyApp/1.0",
    "Samsung Galaxy S21",
)
fmt.Println("Access Token:", companyResp.AccessToken)
```

---

### Verifikasi JWT

```go
// Token user biasa
claims, err := client.VerifyAndParseClaims(accessToken)
if err != nil {
    log.Fatal("Token tidak valid:", err)
}
fmt.Println("User ID:", claims.UserID)
fmt.Println("Company ID:", claims.CompanyID)
fmt.Println("Role:", claims.Role)
fmt.Println("Platform:", claims.Platform)
fmt.Println("Token Type:", claims.TokenType)

// Token S2S (service-to-service)
s2sClaims, err := client.VerifyAndParseS2SClaims(s2sToken)
if err != nil {
    log.Fatal("S2S token tidak valid:", err)
}
fmt.Println("Service Name:", s2sClaims.ServiceName)
```

---

### Refresh Token

```go
resp, err := client.RefreshAccessToken(refreshToken)
if err != nil {
    log.Fatal(err)
}
fmt.Println("New Access Token:", resp.AccessToken)
```

---

### One-Time-Code (OTC) — Cross-platform Navigation

```go
// Generate OTC
otcResp, err := client.GenerateOTCToken(accessToken, authclient.PlatformOldWeb, "", "")
// Dengan company context (multicompany new_web):
// otcResp, err := client.GenerateOTCToken(accessToken, authclient.PlatformOldWeb, "company-id", "employee")

// Exchange OTC untuk access token
exchangeResp, err := client.ExchangeOTCForToken(otcResp.OTCToken)
fmt.Println("Access Token:", exchangeResp.AccessToken)
```

---

### Logout

```go
// Logout sesi ini
_, err := client.Logout(accessToken)

// Logout semua perangkat
_, err = client.LogoutAllDevices(accessToken)
```

---

### Informasi User

```go
userInfo, err := client.GetUserInfo(accessToken)
fmt.Println("Nama:", userInfo.Name)

// Daftar company user (mobile)
companies, err := client.GetUserCompanies(accessToken)

// Semua company (lintas platform)
allCompanies, err := client.GetUserAllCompanies(accessToken)

// Company untuk web
webCompanies, err := client.GetUserWebCompanies(accessToken)
```

---

### Daftarkan Company Baru

```go
phone := "081234567890"
resp, err := client.RegisterCompany(accessToken, "PT Contoh", "admin@contoh.com", &phone)
// phone bisa nil jika tidak ada
```

---

### Manajemen User-Company

```go
// Tambah user ke company sebagai employee (butuh role admin/subadmin/owner)
_, err := client.LinkUserToCompanyAsEmployee(accessToken, "company-uuid", "user-uuid")

// Tambah user ke company sebagai subadmin (butuh role admin)
_, err = client.LinkUserToCompanyAsSubAdmin(accessToken, "company-uuid", "user-uuid")

// Tambah user ke company sebagai owner
_, err = client.LinkUserToCompanyAsOwner(accessToken, "company-uuid", "user-uuid")

// Hapus user dari company (role employee)
_, err = client.UnlinkUserFromCompanyAsEmployee(accessToken, "company-uuid", "user-uuid")
```

---

### OTP

```go
// Request OTP (membutuhkan X-Company-ID)
_, err := client.OTPRequest(
    accessToken, "company-uuid",
    authclient.OTPVerifyTypeLogin,
    authclient.OTPVerifyModePhone,
)

// Verifikasi OTP
resp, err := client.OTPVerify(
    accessToken, "company-uuid", "123456",
    authclient.OTPVerifyTypeLogin,
    authclient.OTPVerifyModePhone,
)

// Verifikasi email
_, err = client.OTPRequestEmailVerification(accessToken)
resp, err = client.OTPVerifyEmail(accessToken, "123456")

// Verifikasi phone
_, err = client.OTPRequestPhoneVerification(accessToken)
resp, err = client.OTPVerifyPhone(accessToken, "123456")
```

---

### Service-to-Service (S2S)

```go
// Issue S2S token
s2sResp, err := client.S2SIssueToken("my-service")
s2sToken := s2sResp.AccessToken

// Reset password via email
_, err = client.S2SRequestEmailResetPassword(s2sToken, "user@example.com", "https://app.example.com")

// Reset password via OTP phone
_, err = client.S2SRequestPhoneOTPResetPassword(s2sToken, "62", "81234567890")

// Eksekusi reset password
_, err = client.S2SResetPassword(s2sToken, "email-reset-token", "", "newpassword123")

// Register company + user admin sekaligus
_, err = client.S2SRegisterCompanyAndUserAdmin(s2sToken, map[string]any{
    "company_name": "PT Baru",
    "email":        "admin@baru.com",
    "password":     "secret123",
})

// Register user saja
_, err = client.S2SRegisterUser(s2sToken, map[string]any{
    "email":    "karyawan@baru.com",
    "password": "secret123",
    "name":     "Budi",
})
```

---

## gRPC Client

gRPC client digunakan untuk query internal antar service (server-to-server). Koneksi dibuat secara lazy (pertama kali method gRPC dipanggil) dan di-reuse untuk semua panggilan berikutnya.

### Inisialisasi dengan gRPC

```go
client := authclient.NewFioAuthClient(
    "http://localhost:8080",  // HTTP base URL
    "localhost:50051",        // gRPC server address
    "my-s2s-api-key",         // API key (dikirim sebagai metadata "authorization")
    30*time.Second,
)

// Untuk environment development (tanpa TLS):
client.WithGRPCInsecure()

// Pastikan koneksi ditutup saat program selesai:
defer client.Close()
```

> **Produksi**: TLS diaktifkan secara default (TLS 1.2+).  
> **Development/local**: Panggil `WithGRPCInsecure()` sebelum request gRPC pertama.  
> **API key**: Dikirim sebagai metadata header `authorization` pada setiap panggilan gRPC. Kosongkan jika tidak dipakai.

---

### GrpcCheckUser — Cek keberadaan user

Memeriksa apakah user dengan ID tertentu ada, dan mengembalikan data profil dasar jika ditemukan.

```go
ctx := context.Background()

result, err := client.GrpcCheckUser(ctx, "user-uuid")
if err != nil {
    log.Fatal(err)
}

if result.Found {
    fmt.Println("Nama:", result.User.Name)
    fmt.Println("Email:", result.User.Email)
    fmt.Println("Status:", result.User.Status)
} else {
    fmt.Println("User tidak ditemukan")
}
```

**Tipe kembalian `GrpcCheckUserResult`:**

| Field   | Tipe            | Keterangan                        |
|---------|-----------------|-----------------------------------|
| `Found` | `bool`          | `true` jika user ditemukan        |
| `User`  | `*GrpcUserBasic`| `nil` jika `Found` adalah `false` |

**Tipe `GrpcUserBasic`:**

| Field       | Tipe     |
|-------------|----------|
| `ID`        | `string` |
| `Name`      | `string` |
| `Email`     | `string` |
| `PhoneCode` | `string` |
| `Phone`     | `string` |
| `Status`    | `string` |

---

### GrpcCheckUserCompanyRelations — Cek semua relasi user di company

Mengembalikan semua role yang dimiliki `userID` di dalam `companyID`.

```go
result, err := client.GrpcCheckUserCompanyRelations(ctx, "user-uuid", "company-uuid")
if err != nil {
    log.Fatal(err)
}

if result.Found {
    for _, rel := range result.Relations {
        fmt.Println("Role:", rel.Role, "| Since:", rel.CreatedAt)
    }
} else {
    fmt.Println("User tidak terdaftar di company ini")
}
```

**Tipe kembalian `GrpcCheckUserCompanyRelationsResult`:**

| Field       | Tipe                        | Keterangan                                     |
|-------------|-----------------------------|------------------------------------------------|
| `Found`     | `bool`                      | `false` jika user tidak punya relasi di company |
| `Relations` | `[]GrpcUserCompanyRelation` | Daftar relasi (bisa lebih dari satu role)       |

---

### GrpcCheckUserCompanyRole — Cek role spesifik user di company

Memeriksa apakah `userID` memiliki role tertentu di `companyID`.

```go
result, err := client.GrpcCheckUserCompanyRole(ctx, "user-uuid", "company-uuid", authclient.RoleEmployee)
if err != nil {
    log.Fatal(err)
}

if result.Found {
    fmt.Println("User adalah employee sejak:", result.Relation.CreatedAt)
} else {
    fmt.Println("User bukan employee di company ini")
}
```

**Nilai `Role` yang valid:** `RoleEmployee`, `RoleOwner`, `RoleSubadmin`, `RoleAdmin`

**Tipe kembalian `GrpcCheckUserCompanyRoleResult`:**

| Field      | Tipe                      | Keterangan                           |
|------------|---------------------------|--------------------------------------|
| `Found`    | `bool`                    | `true` jika user punya role tersebut |
| `Relation` | `*GrpcUserCompanyRelation`| `nil` jika `Found` adalah `false`    |

---

### GrpcGetUserAllRelations — Ambil semua relasi user lintas company

Mengembalikan semua relasi company yang dimiliki `userID` di semua company.

```go
result, err := client.GrpcGetUserAllRelations(ctx, "user-uuid")
if err != nil {
    log.Fatal(err)
}

if result.Found {
    for _, rel := range result.Relations {
        fmt.Printf("Company: %s | Role: %s\n", rel.CompanyID, rel.Role)
    }
}
```

**Tipe kembalian `GrpcGetUserAllRelationsResult`:**

| Field       | Tipe                        | Keterangan                               |
|-------------|-----------------------------|------------------------------------------|
| `Found`     | `bool`                      | `false` jika user tidak punya relasi     |
| `Relations` | `[]GrpcUserCompanyRelation` | Semua relasi user di semua company        |

---

**Tipe `GrpcUserCompanyRelation`:**

| Field       | Tipe        | Keterangan                          |
|-------------|-------------|-------------------------------------|
| `UserID`    | `string`    |                                     |
| `CompanyID` | `string`    |                                     |
| `Role`      | `Role`      | Salah satu konstanta `Role`         |
| `CreatedAt` | `time.Time` | Dikonversi dari unix timestamp      |

---

## Konstanta

| Tipe                | Nilai                                                                                      |
|---------------------|--------------------------------------------------------------------------------------------|
| `Platform`          | `PlatformNewWeb`, `PlatformOldWeb`, `PlatformMobile`, `PlatformPayment`                   |
| `TokenType`         | `TokenTypeAccess`, `TokenTypeRefresh`, `TokenTypeIdentityAccess`, `TokenTypeIdentityRefresh`, `TokenTypeOTC`, `TokenTypeS2SAccess` |
| `OTPVerifyType`     | `OTPVerifyTypeRegister`, `OTPVerifyTypeLogin`, `OTPVerifyTypeResetPassword`, `OTPVerifyTypeEmail`, `OTPVerifyTypePhone`, `OTPVerifyTypeChangeDevice`, `OTPVerifyTypeOther` |
| `OTPVerifyMode`     | `OTPVerifyModePhone`, `OTPVerifyModeEmail`                                                 |
| `Role`              | `RoleEmployee`, `RoleOwner`, `RoleSubadmin`, `RoleAdmin`                                   |
| `BackendModeEnum`   | `BackendModeNewWeb`, `BackendModeOldWeb`                                                   |

---

## Catatan

- JWKS di-cache secara otomatis (default 5 menit). Key rotation ditangani dengan invalidasi cache dan retry otomatis.
- Error HTTP >= 400 dikembalikan sebagai `error` dengan pesan dari server.
- `WebLogin` mengembalikan `"invalid_credentials"` sebagai string error untuk kemudahan assertion.
- Fungsi yang membutuhkan company context (`OTPRequest`, `LinkUserToCompany`, dll.) akan menyertakan header `X-Company-ID` secara otomatis.
- Koneksi gRPC dibuat secara lazy dan di-reuse; panggil `Close()` saat client tidak lagi dibutuhkan.
- `WithGRPCInsecure()` harus dipanggil **sebelum** request gRPC pertama karena koneksi hanya dibuat sekali.
