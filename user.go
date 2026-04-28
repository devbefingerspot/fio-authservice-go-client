package authclient

import "net/http"

// GetUserInfo — GET /api/v1/user/me
func (c *FioAuthClient) GetUserInfo(accessToken string) (*UserInfoResponse, error) {
	var out UserInfoResponse
	_, err := c.doJSON(http.MethodGet, "/api/v1/user/me", nil, bearerHeader(accessToken), &out)
	return &out, err
}

// GetUserCompanies — GET /api/v1/user/companies
//
// Returns the list of companies for a mobile user.
func (c *FioAuthClient) GetUserCompanies(accessToken string) ([]CompanyList, error) {
	var out []CompanyList
	_, err := c.doJSON(http.MethodGet, "/api/v1/user/companies", nil, bearerHeader(accessToken), &out)
	return out, err
}

// GetUserAllCompanies — GET /api/v1/user/all-companies
//
// Returns all companies linked to the user across all platforms.
func (c *FioAuthClient) GetUserAllCompanies(accessToken string) ([]CompanyList, error) {
	var out []CompanyList
	_, err := c.doJSON(http.MethodGet, "/api/v1/user/all-companies", nil, bearerHeader(accessToken), &out)
	return out, err
}

// GetUserWebCompanies — GET /api/v1/user/web-companies
func (c *FioAuthClient) GetUserWebCompanies(accessToken string) ([]CompanyList, error) {
	var out []CompanyList
	_, err := c.doJSON(http.MethodGet, "/api/v1/user/web-companies", nil, bearerHeader(accessToken), &out)
	return out, err
}

// RegisterCompany — POST /api/v1/company/register
//
// Creates a new company and immediately links the current user as admin.
// Does not require a company context header.
func (c *FioAuthClient) RegisterCompany(accessToken, name, email string, phone *string) (*RegisterCompanyResponse, error) {
	body := map[string]any{
		"name":  name,
		"email": email,
	}
	if phone != nil {
		body["phone"] = *phone
	}
	var out RegisterCompanyResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/company/register", body, bearerHeader(accessToken), &out)
	return &out, err
}

// UpdateUserProfile — POST /api/v1/user/update
//
// Updates the authenticated user's name and/or photo URL.
// Accepts any token type including identity token.
// At least one field must be non-nil.
func (c *FioAuthClient) UpdateUserProfile(accessToken string, name, photoURL *string) (*S2SMessageResponse, error) {
	body := map[string]any{}
	if name != nil {
		body["name"] = *name
	}
	if photoURL != nil {
		body["photo_url"] = *photoURL
	}
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user/update", body, bearerHeader(accessToken), &out)
	return &out, err
}

// ChangeUserEmail — POST /api/v1/user/change-email
//
// Changes the authenticated user's email address.
// An OTP must be requested first via RequestEmailOTP with verify_type "change_email".
// The OTP is sent to the user's current email.
func (c *FioAuthClient) ChangeUserEmail(accessToken, newEmail, otpCode string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user/change-email", map[string]any{
		"new_email": newEmail,
		"otp_code":  otpCode,
	}, bearerHeader(accessToken), &out)
	return &out, err
}

// ChangeUserPhone — POST /api/v1/user/change-phone
//
// Changes the authenticated user's phone number.
// An OTP must be requested first via RequestEmailOTP with verify_type "change_phone".
// The OTP is sent to the user's current email.
func (c *FioAuthClient) ChangeUserPhone(accessToken, newPhoneCode, newPhone, otpCode string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user/change-phone", map[string]any{
		"new_phone_code": newPhoneCode,
		"new_phone":      newPhone,
		"otp_code":       otpCode,
	}, bearerHeader(accessToken), &out)
	return &out, err
}
