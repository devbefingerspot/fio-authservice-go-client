package authclient

import (
	"fmt"
	"net/http"
)

// WebLogin — POST /api/v1/login/web
//
// On credential failure (401/404) the error is "invalid_credentials" for
// easy assertion by callers. Use resp.IsRedirect() to detect platform mismatch.
func (c *FioAuthClient) WebLogin(email, password string, platform Platform) (*WebLoginResponse, error) {
	var out WebLoginResponse
	status, err := c.doJSON(http.MethodPost, "/api/v1/login/web", map[string]any{
		"email":    email,
		"password": password,
		"platform": platform,
	}, nil, &out)
	if err != nil {
		if status == http.StatusUnauthorized || status == http.StatusNotFound {
			return nil, fmt.Errorf("invalid_credentials")
		}
		return nil, err
	}
	return &out, nil
}

// RefreshAccessToken — POST /api/v1/auth/refresh
func (c *FioAuthClient) RefreshAccessToken(refreshToken string) (*RefreshTokenResponse, error) {
	var out RefreshTokenResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/auth/refresh", map[string]any{
		"refresh_token": refreshToken,
	}, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GenerateOTCToken — POST /api/v1/auth/otc/generate
//
// Generates a One-Time-Code token for cross-platform navigation.
// companyID and role are optional: when both are provided the OTC carries that
// company context (for multicompany new_web); otherwise the OTC inherits the
// company context of the current access token.
func (c *FioAuthClient) GenerateOTCToken(accessToken string, targetPlatform Platform, companyID, role string) (*GenerateOTCResponse, error) {
	body := map[string]any{
		"target_platform": targetPlatform,
	}
	if companyID != "" && role != "" {
		body["company_id"] = companyID
		body["role"] = role
	}
	var out GenerateOTCResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/auth/otc/generate", body, bearerHeader(accessToken), &out)
	return &out, err
}

// ExchangeOTCForToken — POST /api/v1/auth/otc/exchange
func (c *FioAuthClient) ExchangeOTCForToken(otcToken string) (*ExchangeOTCResponse, error) {
	var out ExchangeOTCResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/auth/otc/exchange", nil, bearerHeader(otcToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Logout — POST /api/v1/auth/logout
func (c *FioAuthClient) Logout(accessToken string) (*LogoutResponse, error) {
	var out LogoutResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/auth/logout", nil, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// LogoutAllDevices — POST /api/v1/auth/logout-all
func (c *FioAuthClient) LogoutAllDevices(accessToken string) (*LogoutResponse, error) {
	var out LogoutResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/auth/logout-all", nil, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
