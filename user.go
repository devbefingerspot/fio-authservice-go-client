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
