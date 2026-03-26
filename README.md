
# fio-auth-service-go-client

Go client library for Fingerspot Auth Service. Provides login functions, JWT verification, OTP, S2S token, user-company management, and a gRPC client for internal queries.

---


## Installation

```bash
go get github.com/devbefingerspot/fio-authservice-go-client
```

---


## Client Initialization


```go
import authclient "github.com/devbefingerspot/fio-authservice-go-client"

client := authclient.NewFioAuthClient(
    "http://localhost:8080",       // base URL of auth service (HTTP)
    "auth-grpc.example.com:50051", // gRPC server base URL (leave empty to use the same host)
    "my-api-key",                  // API key for gRPC (leave empty if not used)
    "my-s2s-key",                  // S2S pre-shared key (required for S2S token)
    30*time.Second,                 // HTTP timeout
    // optional: JWKS cache TTL (default 5 minutes)
    // 10*time.Minute,
)
defer client.Close() // close gRPC connection when done
```

### Additional Options

```go
// Disable TLS on gRPC connection (for development/local only)
client.WithGRPCInsecure()
```

---


## Usage Examples

### Health Check

```go
resp, err := client.HealthCheck()
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Message)
```

---


### Web Login

```go
resp, err := client.WebLogin("user@example.com", "password123", authclient.PlatformNewWeb)
if err != nil {
    if err.Error() == "invalid_credentials" {
        log.Println("Email or password is incorrect")
    }
    log.Fatal(err)
}

if resp.IsRedirect() {
    // Platform mismatch — user is redirected to another platform
    fmt.Println("Redirect to:", resp.RedirectPlatform)
    fmt.Println("OTC Token:", *resp.OTCToken)
} else {
    fmt.Println("Access Token:", resp.AccessToken)
    fmt.Println("Refresh Token:", resp.RefreshToken)
}
```

---


### Mobile Login

```go

// Login with email
resp, err := client.MobileLogin("password123", "user@example.com", "", "")

// Login with phone number
resp, err := client.MobileLogin("password123", "", "81234567890", "62")

if err != nil {
    log.Fatal(err)
}
fmt.Println("Identity Access Token:", resp.IdentityAccessToken)

// Exchange identity token for company-scoped token
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


### JWT Verification

```go
// Regular user token
claims, err := client.VerifyAndParseClaims(accessToken)
if err != nil {
    log.Fatal("Invalid token:", err)
}
fmt.Println("User ID:", claims.UserID)
fmt.Println("Company ID:", claims.CompanyID)
fmt.Println("Role:", claims.Role)
fmt.Println("Platform:", claims.Platform)
fmt.Println("Token Type:", claims.TokenType)

// S2S (service-to-service) token
s2sClaims, err := client.VerifyAndParseS2SClaims(s2sToken)
if err != nil {
    log.Fatal("Invalid S2S token:", err)
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
// With company context (multicompany new_web):
// otcResp, err := client.GenerateOTCToken(accessToken, authclient.PlatformOldWeb, "company-id", "employee")

// Exchange OTC for access token
exchangeResp, err := client.ExchangeOTCForToken(otcResp.OTCToken)
fmt.Println("Access Token:", exchangeResp.AccessToken)
```

---


### Logout

```go
// Logout this session
_, err := client.Logout(accessToken)

// Logout all devices
_, err = client.LogoutAllDevices(accessToken)
```

---


### User Information

```go
userInfo, err := client.GetUserInfo(accessToken)
fmt.Println("Name:", userInfo.Name)

// List of user's companies (mobile)
companies, err := client.GetUserCompanies(accessToken)

// All companies (cross-platform)
allCompanies, err := client.GetUserAllCompanies(accessToken)

// Companies for web
webCompanies, err := client.GetUserWebCompanies(accessToken)
```

---


### Register a New Company

```go
phone := "081234567890"
resp, err := client.RegisterCompany(accessToken, "PT Example", "admin@example.com", &phone)
// phone can be nil if not provided
```

---


### User-Company Management

```go
// Add user to company as employee (requires admin/subadmin/owner role)
_, err := client.LinkUserToCompanyAsEmployee(accessToken, "company-uuid", "user-uuid")

// Add user to company as subadmin (requires admin role)
_, err = client.LinkUserToCompanyAsSubAdmin(accessToken, "company-uuid", "user-uuid")

// Add user to company as owner
_, err = client.LinkUserToCompanyAsOwner(accessToken, "company-uuid", "user-uuid")

// Remove user from company (employee role)
_, err = client.UnlinkUserFromCompanyAsEmployee(accessToken, "company-uuid", "user-uuid")
```

---


### OTP

```go
// Request OTP (requires X-Company-ID)
_, err := client.OTPRequest(
    accessToken, "company-uuid",
    authclient.OTPVerifyTypeLogin,
    authclient.OTPVerifyModePhone,
)

// Verify OTP
resp, err := client.OTPVerify(
    accessToken, "company-uuid", "123456",
    authclient.OTPVerifyTypeLogin,
    authclient.OTPVerifyModePhone,
)

// Email verification
_, err = client.OTPRequestEmailVerification(accessToken)
resp, err = client.OTPVerifyEmail(accessToken, "123456")

// Phone verification
_, err = client.OTPRequestPhoneVerification(accessToken)
resp, err = client.OTPVerifyPhone(accessToken, "123456")
```

---


### Service-to-Service (S2S)

```go
// Issue S2S token (requires s2sKey set in client constructor)
s2sResp, err := client.S2SIssueToken("my-service")
s2sToken := s2sResp.AccessToken

// Reset password via email
_, err = client.S2SRequestEmailResetPassword(s2sToken, "user@example.com", "https://app.example.com")

// Reset password via OTP phone
_, err = client.S2SRequestPhoneOTPResetPassword(s2sToken, "62", "81234567890")

// Execute password reset
_, err = client.S2SResetPassword(s2sToken, "email-reset-token", "", "newpassword123")

// Register company + admin user at once
_, err = client.S2SRegisterCompanyAndUserAdmin(s2sToken, map[string]any{
    "company_name": "PT New",
    "email":        "admin@new.com",
    "password":     "secret123",
})

// Register user only
_, err = client.S2SRegisterUser(s2sToken, map[string]any{
    "email":    "employee@new.com",
    "password": "secret123",
    "name":     "Budi",
})
```

> **S2S Security:**
>
> - Untuk generate S2S token, client harus mengisi parameter `s2sKey` pada constructor `NewFioAuthClient`.
> - Library otomatis mengirim header `X-S2S-Authorization` saat request S2S token.
> - Jika key salah atau tidak diisi, request akan gagal (401 Unauthorized).

---

## gRPC Client


gRPC client is used for internal queries between services (server-to-server). The connection is created lazily (on the first gRPC method call) and reused for all subsequent calls.


### Initialization with gRPC

```go
client := authclient.NewFioAuthClient(
    "http://localhost:8080",  // HTTP base URL
    "localhost:50051",        // gRPC server address
    "my-s2s-api-key",         // API key (sent as "authorization" metadata)
    30*time.Second,
)

// For development environment (without TLS):
client.WithGRPCInsecure()

// Make sure to close the connection when the program ends:
defer client.Close()
```

> **Production**: TLS is enabled by default (TLS 1.2+).  
> **Development/local**: Call `WithGRPCInsecure()` before the first gRPC request.  
> **API key**: Sent as `authorization` metadata header on every gRPC call. Leave empty if not used.

---


### GrpcCheckUser — Check user existence

Checks if a user with the given ID exists, and returns basic profile data if found.

```go
ctx := context.Background()

result, err := client.GrpcCheckUser(ctx, "user-uuid")
if err != nil {
    log.Fatal(err)
}

if result.Found {
    fmt.Println("Name:", result.User.Name)
    fmt.Println("Email:", result.User.Email)
    fmt.Println("Status:", result.User.Status)
} else {
    fmt.Println("User tidak ditemukan")
}
```


**Return type `GrpcCheckUserResult`:**

| Field   | Type             | Description                          |
|---------|------------------|--------------------------------------|
| `Found` | `bool`           | `true` if user is found              |
| `User`  | `*GrpcUserBasic` | `nil` if `Found` is `false`          |

**Type `GrpcUserBasic`:**

| Field             | Type    | Description                                         |
|-------------------|---------|-----------------------------------------------------|
| `ID`              | `string`|                                                     |
| `Name`            | `string`|                                                     |
| `Email`           | `string`|                                                     |
| `PhoneCode`       | `string`|                                                     |
| `Phone`           | `string`|                                                     |
| `Status`          | `string`|                                                     |
| `EmailVerifiedAt` | `int64` | Unix timestamp; `0` if not verified                 |
| `PhoneVerifiedAt` | `int64` | Unix timestamp; `0` if not verified                 |

---


### GrpcCheckUserCompanyRelations — Check all user relations in a company

Returns all roles owned by `userID` in `companyID`.

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
    fmt.Println("User is not registered in this company")
}
```


**Return type `GrpcCheckUserCompanyRelationsResult`:**

| Field       | Type                        | Description                                     |
|-------------|-----------------------------|-------------------------------------------------|
| `Found`     | `bool`                      | `false` if user has no relation in the company  |
| `Relations` | `[]GrpcUserCompanyRelation` | List of relations (can have more than one role) |

---


### GrpcCheckUserCompanyRole — Check specific user role in a company

Checks if `userID` has a specific role in `companyID`.

```go
result, err := client.GrpcCheckUserCompanyRole(ctx, "user-uuid", "company-uuid", authclient.RoleEmployee)
if err != nil {
    log.Fatal(err)
}

if result.Found {
    fmt.Println("User is an employee since:", result.Relation.CreatedAt)
} else {
    fmt.Println("User is not an employee in this company")
}
```


**Valid `Role` values:** `RoleEmployee`, `RoleOwner`, `RoleSubadmin`, `RoleAdmin`

**Return type `GrpcCheckUserCompanyRoleResult`:**

| Field      | Type                      | Description                           |
|------------|---------------------------|---------------------------------------|
| `Found`    | `bool`                    | `true` if user has the role           |
| `Relation` | `*GrpcUserCompanyRelation`| `nil` if `Found` is `false`           |

---


### GrpcGetUserAllRelations — Get all user relations across companies

Returns all company relations owned by `userID` in all companies.

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


**Return type `GrpcGetUserAllRelationsResult`:**

| Field       | Type                        | Description                               |
|-------------|-----------------------------|-------------------------------------------|
| `Found`     | `bool`                      | `false` if user has no relations          |
| `Relations` | `[]GrpcUserCompanyRelation` | All user relations in all companies       |

---


**Type `GrpcUserCompanyRelation`:**

| Field       | Type        | Description                          |
|-------------|-------------|--------------------------------------|
| `UserID`    | `string`    |                                      |
| `CompanyID` | `string`    |                                      |
| `Role`      | `Role`      | One of the `Role` constants          |
| `CreatedAt` | `time.Time` | Converted from unix timestamp        |

---


### GrpcGetCompanyWithEndpoint — Get company info with endpoint

Returns company and related endpoint data based on `companyID`.

```go
result, err := client.GrpcGetCompanyWithEndpoint(ctx, "company-uuid")
if err != nil {
    log.Fatal(err)
}

if result.Found {
    fmt.Println("Company:", result.Company.Name)
    fmt.Println("Email:", result.Company.Email)
    fmt.Println("Device Policy:", result.Company.DeviceLoginPolicy)
    fmt.Println("Max Devices:", result.Company.MaxDevices)

    if result.Endpoint != nil {
        fmt.Println("Backend Mode:", result.Endpoint.BackendMode)
        fmt.Println("Base URL:", result.Endpoint.BaseURL)
        fmt.Println("DB Driver:", result.Endpoint.DBDriver) // kosong jika tidak diset
        fmt.Println("DB DSN:", result.Endpoint.DBDSN)       // kosong jika tidak diset
    }
} else {
    fmt.Println("Company not found")
}
```


**Return type `GrpcGetCompanyWithEndpointResult`:**

| Field      | Type                | Description                                 |
|------------|---------------------|---------------------------------------------|
| `Found`    | `bool`              | `false` if company not found                |
| `Company`  | `*GrpcCompanyInfo`  | `nil` if `Found` is `false`                 |
| `Endpoint` | `*GrpcEndpointInfo` | `nil` if company has no endpoint            |

**Type `GrpcCompanyInfo`:**

| Field               | Type     | Description                                 |
|---------------------|----------|---------------------------------------------|
| `ID`                | `string` |                                             |
| `Name`              | `string` |                                             |
| `Email`             | `string` |                                             |
| `Phone`             | `string` | Empty if not set                            |
| `DueDate`           | `int64`  | Unix timestamp; `0` if not set              |
| `EndpointID`        | `string` |                                             |
| `DeviceLoginPolicy` | `string` | `fixed_device` or `trusted_device_rotate`   |
| `MaxDevices`        | `int32`  |                                             |

**Type `GrpcEndpointInfo`:**

| Field         | Type     | Description                        |
|---------------|----------|------------------------------------|
| `ID`          | `string` |                                    |
| `BackendMode` | `string` | `new_web` or `old_web`             |
| `BaseURL`     | `string` |                                    |
| `DBDriver`    | `string` | Empty if not set                   |
| `DBDSN`       | `string` | Empty if not set                   |

---


## Constants

| Type                | Values                                                                                      |
|---------------------|--------------------------------------------------------------------------------------------|
| `Platform`          | `PlatformNewWeb`, `PlatformOldWeb`, `PlatformMobile`, `PlatformPayment`                    |
| `TokenType`         | `TokenTypeAccess`, `TokenTypeRefresh`, `TokenTypeIdentityAccess`, `TokenTypeIdentityRefresh`, `TokenTypeOTC`, `TokenTypeS2SAccess` |
| `OTPVerifyType`     | `OTPVerifyTypeRegister`, `OTPVerifyTypeLogin`, `OTPVerifyTypeResetPassword`, `OTPVerifyTypeEmail`, `OTPVerifyTypePhone`, `OTPVerifyTypeChangeDevice`, `OTPVerifyTypeOther` |
| `OTPVerifyMode`     | `OTPVerifyModePhone`, `OTPVerifyModeEmail`                                                 |
| `Role`              | `RoleEmployee`, `RoleOwner`, `RoleSubadmin`, `RoleAdmin`                                   |
| `BackendModeEnum`   | `BackendModeNewWeb`, `BackendModeOldWeb`                                                   |

---


## Notes

- JWKS is automatically cached (default 5 minutes). Key rotation is handled with cache invalidation and automatic retry.
- HTTP errors (>= 400) are returned as `error` with the server message.
- `WebLogin` returns `"invalid_credentials"` as a string error for easier assertion.
- Functions requiring company context (`OTPRequest`, `LinkUserToCompany`, etc.) will automatically include the `X-Company-ID` header.
- gRPC connection is created lazily and reused; call `Close()` when the client is no longer needed.
- `WithGRPCInsecure()` must be called **before** the first gRPC request because the connection is created only once.
