# fio-auth-service-go-client

Go client library untuk Fingerspot Auth Service. Menyediakan fungsi login, verifikasi JWT, OTP, S2S token, dan manajemen user-company.

---

## Instalasi

```bash
go get local/fio-auth-service-client
```

> Sesuaikan module path dengan konfigurasi `go.mod` di project kamu.

---

## Inisialisasi Client

```go
import authclient "local/fio-auth-service-client"

client := authclient.NewFioAuthClient(
    "http://localhost:8080", // base URL auth service
    30*time.Second,          // HTTP timeout
    // optional: cache TTL untuk JWKS (default 5 menit)
    // 10*time.Minute,
)
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

### Manajemen User-Company

```go
// Tambah user ke company
_, err := client.LinkUserToCompanyAsEmployee(accessToken, companyID, userID)
_, err = client.LinkUserToCompanyAsSubAdmin(accessToken, companyID, userID)
_, err = client.LinkUserToCompanyAsOwner(accessToken, companyID, userID)

// Hapus user dari company
_, err = client.UnlinkUserFromCompanyAsEmployee(accessToken, companyID, userID)
_, err = client.UnlinkUserFromCompanyAsSubAdmin(accessToken, companyID, userID)
```

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
