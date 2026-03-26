package authclient

import "net/http"

// LinkUserToCompanyAsEmployee — POST /api/v1/user-company/link/employee
//
// Adds a user to the company with the employee role.
// Requires admin, subadmin, or owner access.
func (c *FioAuthClient) LinkUserToCompanyAsEmployee(accessToken, companyID, userID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user-company/link/employee", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// LinkUserToCompanyAsSubAdmin — POST /api/v1/user-company/link/subadmin
//
// Adds a user to the company with the subadmin role.
// Requires admin access.
func (c *FioAuthClient) LinkUserToCompanyAsSubAdmin(accessToken, companyID, userID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user-company/link/subadmin", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// LinkUserToCompanyAsOwner — POST /api/v1/user-company/link/owner
//
// Adds a user to the company with the owner role.
func (c *FioAuthClient) LinkUserToCompanyAsOwner(accessToken, companyID, userID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user-company/link/owner", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// UnlinkUserFromCompanyAsEmployee — DELETE /api/v1/user-company/unlink/employee
//
// Removes a user from the company (employee role).
func (c *FioAuthClient) UnlinkUserFromCompanyAsEmployee(accessToken, companyID, userID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodDelete, "/api/v1/user-company/unlink/employee", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// UnlinkUserFromCompanyAsSubAdmin — DELETE /api/v1/user-company/unlink/subadmin
//
// Removes a user from the company (subadmin role).
func (c *FioAuthClient) UnlinkUserFromCompanyAsSubAdmin(accessToken, companyID, userID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodDelete, "/api/v1/user-company/unlink/subadmin", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// UnlinkUserFromCompanyAsOwner — DELETE /api/v1/user-company/unlink/owner
//
// Removes a user from the company (owner role).
func (c *FioAuthClient) UnlinkUserFromCompanyAsOwner(accessToken, companyID, userID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodDelete, "/api/v1/user-company/unlink/owner", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// ChangeEmployeeDevice — POST /api/v1/user-device/change
//
// Updates the registered device for an employee.
// Requires admin, subadmin, or owner access.
func (c *FioAuthClient) ChangeEmployeeDevice(
	accessToken, companyID, userID string,
	deviceUUID, fcmToken, userAgent, detail string,
) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user-device/change", map[string]any{
		"user_id":                  userID,
		"device_unique_identifier": deviceUUID,
		"fcm_token":                fcmToken,
		"user_agent":               userAgent,
		"detail":                   detail,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}

// ChangeCompanyEndpoint — POST /api/v1/company/change-endpoint
//
// Changes the backend mode of the company. Valid values: "new_web" or "old_web".
func (c *FioAuthClient) ChangeCompanyEndpoint(accessToken, companyID, backendMode string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/company/change-endpoint", map[string]any{
		"backend_mode": backendMode,
	}, companyContextHeaders(accessToken, companyID), &out)
	return &out, err
}
