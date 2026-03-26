package authclient

import "net/http"

// MobileLogin — POST /api/v1/login/mobile
//
// Provide either email OR phone+phoneCode (leave the unused fields empty).
func (c *FioAuthClient) MobileLogin(password, email, phone, phoneCode string) (*MobileLoginResponse, error) {
	body := map[string]any{"password": password}
	if email != "" {
		body["email"] = email
	}
	if phone != "" {
		body["phone"] = phone
		body["phone_code"] = phoneCode
	}

	var out MobileLoginResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/login/mobile", body, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// MobileIssueCompanyToken — POST /api/v1/mobile/issue-company-token
//
// Exchanges an identity access token for a company-scoped access token on mobile.
func (c *FioAuthClient) MobileIssueCompanyToken(
	identityAccessToken, companyID, role, deviceUUID, fcmToken, userAgent, detail string,
) (*MobileIssueCompanyTokenResponse, error) {
	var out MobileIssueCompanyTokenResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/mobile/issue-company-token", map[string]any{
		"company_id":               companyID,
		"role":                     role,
		"device_unique_identifier": deviceUUID,
		"fcm_token":                fcmToken,
		"user_agent":               userAgent,
		"detail":                   detail,
	}, bearerHeader(identityAccessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
