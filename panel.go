package authclient

import (
	"fmt"
	"net/http"
)

// PanelLogin — POST /api/v1/panel/auth/login
//
// Authenticates a panel user (internal/channels) and returns access and refresh tokens.
func (c *FioAuthClient) PanelLogin(email, password string) (*PanelLoginResponse, error) {
	var out PanelLoginResponse
	status, err := c.doJSON(http.MethodPost, "/api/v1/panel/auth/login", map[string]any{
		"email":    email,
		"password": password,
	}, nil, &out)
	if err != nil {
		if status == http.StatusUnauthorized || status == http.StatusNotFound {
			return nil, fmt.Errorf("invalid_credentials")
		}
		return nil, err
	}
	return &out, nil
}

// PanelRefreshAccessToken — POST /api/v1/panel/auth/refresh
func (c *FioAuthClient) PanelRefreshAccessToken(refreshToken string) (*PanelRefreshTokenResponse, error) {
	var out PanelRefreshTokenResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/panel/auth/refresh", map[string]any{
		"refresh_token": refreshToken,
	}, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// PanelLogout — POST /api/v1/panel/auth/logout
func (c *FioAuthClient) PanelLogout(accessToken string) (*LogoutResponse, error) {
	var out LogoutResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/panel/auth/logout", nil, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// PanelLogoutAllDevices — POST /api/v1/panel/auth/logout-all
func (c *FioAuthClient) PanelLogoutAllDevices(accessToken string) (*LogoutResponse, error) {
	var out LogoutResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/panel/auth/logout-all", nil, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// PanelMe — GET /api/v1/panel/me
//
// Returns the current panel user's profile and roles.
func (c *FioAuthClient) PanelMe(accessToken string) (*PanelMeResponse, error) {
	var out PanelMeResponse
	_, err := c.doJSON(http.MethodGet, "/api/v1/panel/me", nil, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
