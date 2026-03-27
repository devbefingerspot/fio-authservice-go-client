package authclient

// Platform identifies which platform a token was issued for.
type Platform string

const (
	PlatformNewWeb  Platform = "new_web"
	PlatformOldWeb  Platform = "old_web"
	PlatformMobile  Platform = "mobile"
	PlatformPayment Platform = "payment"
)

// TokenType identifies the purpose of a JWT token.
type TokenType string

const (
	TokenTypeAccess          TokenType = "access"
	TokenTypeRefresh         TokenType = "refresh"
	TokenTypeIdentityAccess  TokenType = "identity_access"
	TokenTypeIdentityRefresh TokenType = "identity_refresh"
	TokenTypeOTC             TokenType = "otc"
	TokenTypeS2SAccess       TokenType = "s2s_access"
)

// OTPVerifyType identifies the purpose of an OTP verification.
type OTPVerifyType string

const (
	OTPVerifyTypeRegister      OTPVerifyType = "register"
	OTPVerifyTypeLogin         OTPVerifyType = "login"
	OTPVerifyTypeResetPassword OTPVerifyType = "reset_password"
	OTPVerifyTypeEmail         OTPVerifyType = "email_verification"
	OTPVerifyTypePhone         OTPVerifyType = "phone_verification"
	OTPVerifyTypeChangeDevice  OTPVerifyType = "change_device"
	OTPVerifyTypeOther         OTPVerifyType = "other"
)

// OTPVerifyMode identifies the delivery channel for an OTP.
type OTPVerifyMode string

const (
	OTPVerifyModePhone OTPVerifyMode = "phone"
	OTPVerifyModeEmail OTPVerifyMode = "email"
)

// Role represents a user's role within a company.
type Role string

const (
	RoleEmployee Role = "employee"
	RoleOwner    Role = "owner"
	RoleSubadmin Role = "subadmin"
	RoleAdmin    Role = "admin"
)

// BackendModeEnum identifies the backend mode (endpoint type) of a company.
type BackendModeEnum string

const (
	BackendModeNewWeb BackendModeEnum = "new_web"
	BackendModeOldWeb BackendModeEnum = "old_web"
)

// SubjectType identifies the type of subject in a JWT token.
type SubjectType string

const (
	SubjectTypeAppUser   SubjectType = "app_user"
	SubjectTypePanelUser SubjectType = "panel_user"
)

// PanelRole represents a panel user's role.
type PanelRole string

const (
	PanelRoleInternal PanelRole = "internal"
	PanelRoleChannels PanelRole = "channels"
)

// PanelUserStatus represents the status of a panel user account.
type PanelUserStatus string

const (
	PanelUserStatusActive    PanelUserStatus = "active"
	PanelUserStatusInactive  PanelUserStatus = "inactive"
	PanelUserStatusSuspended PanelUserStatus = "suspended"
)
