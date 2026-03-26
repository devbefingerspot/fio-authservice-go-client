package authclient

import "time"

// ErrorResponse represents any error body returned by the auth-service.
type ErrorResponse struct {
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   string `json:"details,omitempty"`
}

func (e *ErrorResponse) String() string {
	if e.Error != "" {
		return e.Error
	}
	if e.Message != "" {
		return e.Message
	}
	return e.ErrorCode
}

// HealthCheckResponse — GET /api/v1/
type HealthCheckResponse struct {
	Message string `json:"message"`
}

// WebLoginResponse — POST /api/v1/login/web
// Status is "success" or "redirect" (platform mismatch).
type WebLoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	// success fields
	UserID                string `json:"user_id,omitempty"`
	UserName              string `json:"user_name,omitempty"`
	CompanyID             string `json:"company_id,omitempty"`
	CompanyName           string `json:"company_name,omitempty"`
	Role                  string `json:"role,omitempty"`
	OldCompanyID          *int   `json:"old_company_id,omitempty"`
	OldUserID             *int   `json:"old_user_id,omitempty"`
	AccessToken           string `json:"access_token,omitempty"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at,omitempty"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at,omitempty"`

	// redirect fields (platform mismatch)
	OTCToken         *string `json:"otc_token,omitempty"`
	RedirectPlatform string  `json:"redirect_platform,omitempty"`
}

func (r *WebLoginResponse) IsRedirect() bool {
	return r.Status == "redirect"
}

// MobileLoginResponse — POST /api/v1/login/mobile
type MobileLoginResponse struct {
	Message                       string `json:"message"`
	UserID                        string `json:"user_id"`
	UserName                      string `json:"user_name"`
	OldUserID                     *int   `json:"old_user_id,omitempty"`
	IdentityAccessToken           string `json:"identity_access_token"`
	IdentityAccessTokenExpiresAt  int64  `json:"identity_access_token_expires_at"`
	IdentityRefreshToken          string `json:"identity_refresh_token"`
	IdentityRefreshTokenExpiresAt int64  `json:"identity_refresh_token_expires_at"`
}

// MobileIssueCompanyTokenResponse — POST /api/v1/mobile/issue-company-token
type MobileIssueCompanyTokenResponse struct {
	Message               string `json:"message"`
	UserID                string `json:"user_id"`
	CompanyID             string `json:"company_id"`
	CompanyName           string `json:"company_name"`
	Role                  string `json:"role"`
	OldCompanyID          *int   `json:"old_company_id,omitempty"`
	OldUserID             *int   `json:"old_user_id,omitempty"`
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
}

// RefreshTokenResponse — POST /api/v1/auth/refresh
type RefreshTokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
	TokenTypeHint         string `json:"token_type"`
}

// GenerateOTCResponse — POST /api/v1/auth/otc/generate
type GenerateOTCResponse struct {
	Status   string `json:"status"`
	OTCToken string `json:"otc_token"`
	Message  string `json:"message"`
}

// ExchangeOTCResponse — POST /api/v1/auth/otc/exchange
type ExchangeOTCResponse struct {
	Status                string `json:"status"`
	Message               string `json:"message"`
	UserID                string `json:"user_id"`
	CompanyID             string `json:"company_id"`
	Role                  string `json:"role"`
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
}

// LogoutResponse — POST /api/v1/auth/logout and /api/v1/auth/logout-all
type LogoutResponse struct {
	Message string `json:"message"`
}

// S2STokenResponse — POST /api/v1/s2s/token
type S2STokenResponse struct {
	Token string `json:"token"`
}

// OTPRequestResponse — POST /api/v1/otp/request, /otp/email/request, /otp/phone/request
type OTPRequestResponse struct {
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

// OTPVerifyResponse — POST /api/v1/otp/verify, /otp/email/verify, /otp/phone/verify
type OTPVerifyResponse struct {
	Message string `json:"message"`
}

// RegisterCompanyNewCompany holds data for a newly registered company.
type RegisterCompanyNewCompany struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Phone *string `json:"phone,omitempty"`
}

// RegisterCompanyUserCompany holds the user-company relation created on registration.
type RegisterCompanyUserCompany struct {
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id"`
	Role      string `json:"role"`
}

// RegisterCompanyResponse — POST /api/v1/company/register
type RegisterCompanyResponse struct {
	Company     RegisterCompanyNewCompany  `json:"company"`
	UserCompany RegisterCompanyUserCompany `json:"user_company"`
}

// CompanyList represents a company entry in a user's company list.
type CompanyList struct {
	ID          string          `json:"id"`
	OldID       *int            `json:"old_id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Phone       *string         `json:"phone"`
	Role        Role            `json:"role"`
	BackendMode BackendModeEnum `json:"backend_mode"`
	BaseURL     string          `json:"base_url"`
}

// UserInfoUser holds profile details of a user.
type UserInfoUser struct {
	ID              string     `json:"id"`
	OldID           *int       `json:"old_id,omitempty"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	PhoneCode       *string    `json:"phone_code,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	Status          *string    `json:"status,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// UserInfoCompany holds company details within a user info response.
type UserInfoCompany struct {
	ID        string     `json:"id"`
	OldID     *int       `json:"old_id,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone,omitempty"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// UserInfoResponse — GET /api/v1/user/me
type UserInfoResponse struct {
	User    UserInfoUser    `json:"user"`
	Company UserInfoCompany `json:"company"`
	Role    string          `json:"role"`
}

// S2SMessageResponse is a generic message response for S2S endpoints.
type S2SMessageResponse struct {
	Message string `json:"message"`
}

// S2SRegisterUserResponse — POST /api/v1/s2s/user/register
type S2SRegisterUserResponse struct {
	Message       string `json:"message"`
	UserID        string `json:"user_id"`
	UserEmail     string `json:"user_email"`
	UserName      string `json:"user_name"`
	UserPhone     string `json:"user_phone"`
	UserPhoneCode string `json:"user_phone_code"`
	UserOldID     *int   `json:"user_old_id"`
}

// S2SAPIErrorResponse is an error body specific to S2S endpoints.
type S2SAPIErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
