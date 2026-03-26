package authclient

import "net/http"

// S2SIssueToken — POST /api/v1/s2s/token
//
// Issues a service-to-service JWT for the given service name.
func (c *FioAuthClient) S2SIssueToken(serviceName string) (*S2STokenResponse, error) {
	var out S2STokenResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/s2s/token", map[string]any{
		"service_name": serviceName,
	}, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// S2SRequestEmailResetPassword — POST /api/v1/s2s/reset-password/email
func (c *FioAuthClient) S2SRequestEmailResetPassword(s2sToken, email, baseURL string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/s2s/reset-password/email", map[string]any{
		"email":    email,
		"base_url": baseURL,
	}, bearerHeader(s2sToken), &out)
	return &out, err
}

// S2SRequestPhoneOTPResetPassword — POST /api/v1/s2s/reset-password/phone-otp
func (c *FioAuthClient) S2SRequestPhoneOTPResetPassword(s2sToken, phoneCode, phone string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/s2s/reset-password/phone-otp", map[string]any{
		"phone_code": phoneCode,
		"phone":      phone,
	}, bearerHeader(s2sToken), &out)
	return &out, err
}

// S2SResetPassword — POST /api/v1/s2s/reset-password
//
// Provide token (email flow) or otpCode (phone flow); the unused field may be empty.
func (c *FioAuthClient) S2SResetPassword(s2sToken, token, otpCode, newPassword string) (*S2SMessageResponse, error) {
	body := map[string]any{"new_password": newPassword}
	if token != "" {
		body["token"] = token
	}
	if otpCode != "" {
		body["otp_code"] = otpCode
	}
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/s2s/reset-password", body, bearerHeader(s2sToken), &out)
	return &out, err
}

// S2SRegisterCompanyAndUserAdmin — POST /api/v1/s2s/user-company/register
func (c *FioAuthClient) S2SRegisterCompanyAndUserAdmin(s2sToken string, body map[string]any) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/s2s/user-company/register", body, bearerHeader(s2sToken), &out)
	return &out, err
}

// S2SRegisterUser — POST /api/v1/s2s/user/register
func (c *FioAuthClient) S2SRegisterUser(s2sToken string, body map[string]any) (*S2SRegisterUserResponse, error) {
	var out S2SRegisterUserResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/s2s/user/register", body, bearerHeader(s2sToken), &out)
	return &out, err
}
