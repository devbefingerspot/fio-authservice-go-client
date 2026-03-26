package authclient

import "net/http"

// OTPRequest — POST /api/v1/otp/request
//
// Sends an OTP to the user according to verifyType and verifyMode.
// Requires company context (X-Company-ID).
func (c *FioAuthClient) OTPRequest(accessToken, companyID string, verifyType OTPVerifyType, verifyMode OTPVerifyMode) (*OTPRequestResponse, error) {
	var out OTPRequestResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/otp/request", map[string]any{
		"verify_type": verifyType,
		"verify_mode": verifyMode,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// OTPVerify — POST /api/v1/otp/verify
//
// Verifies the OTP code. Requires company context (X-Company-ID).
func (c *FioAuthClient) OTPVerify(accessToken, companyID, code string, verifyType OTPVerifyType, verifyMode OTPVerifyMode) (*OTPVerifyResponse, error) {
	var out OTPVerifyResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/otp/verify", map[string]any{
		"code":        code,
		"verify_type": verifyType,
		"verify_mode": verifyMode,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// OTPRequestEmailVerification — POST /api/v1/otp/email/request
//
// Sends an OTP to the user's email for email verification.
func (c *FioAuthClient) OTPRequestEmailVerification(accessToken string) (*OTPRequestResponse, error) {
	var out OTPRequestResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/otp/email/request", nil, bearerHeader(accessToken), &out)
	return &out, err
}

// OTPVerifyEmail — POST /api/v1/otp/email/verify
func (c *FioAuthClient) OTPVerifyEmail(accessToken, code string) (*OTPVerifyResponse, error) {
	var out OTPVerifyResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/otp/email/verify", map[string]any{
		"code": code,
	}, bearerHeader(accessToken), &out)
	return &out, err
}

// OTPRequestPhoneVerification — POST /api/v1/otp/phone/request
//
// Sends an OTP to the user's phone number for phone verification.
func (c *FioAuthClient) OTPRequestPhoneVerification(accessToken string) (*OTPRequestResponse, error) {
	var out OTPRequestResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/otp/phone/request", nil, bearerHeader(accessToken), &out)
	return &out, err
}

// OTPVerifyPhone — POST /api/v1/otp/phone/verify
func (c *FioAuthClient) OTPVerifyPhone(accessToken, code string) (*OTPVerifyResponse, error) {
	var out OTPVerifyResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/otp/phone/verify", map[string]any{
		"code": code,
	}, bearerHeader(accessToken), &out)
	return &out, err
}
