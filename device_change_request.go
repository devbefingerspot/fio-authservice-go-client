package authclient

import (
	"fmt"
	"net/http"
)

// MyDeviceItem represents a single registered device for a user.
type MyDeviceItem struct {
	DeviceID  string `json:"device_id"`
	FcmToken  string `json:"fcm_token"`
	CompanyID string `json:"company_id"`
	CreatedAt string `json:"created_at"`
}

// MyDevicesResponse — GET /api/v1/user/my-devices
type MyDevicesResponse struct {
	Devices []MyDeviceItem `json:"devices"`
}

// GetMyDevices — GET /api/v1/user/my-devices
//
// Returns all registered mobile devices for the authenticated user.
func (c *FioAuthClient) GetMyDevices(accessToken string) (*MyDevicesResponse, error) {
	var out MyDevicesResponse
	_, err := c.doJSON(http.MethodGet, "/api/v1/user/my-devices", nil, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDeviceChangeRequestResponse — POST /api/v1/device-change-request
type CreateDeviceChangeRequestResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

// CreateDeviceChangeRequest — POST /api/v1/device-change-request
//
// Submits a request to change a mobile device. The user specifies which old device
// to replace and provides the new device details. The request is reviewed by an admin.
func (c *FioAuthClient) CreateDeviceChangeRequest(
	accessToken, companyID, oldDeviceID, newDeviceUUID, newFcmToken, newUserAgent, newDetail string,
) (*CreateDeviceChangeRequestResponse, error) {
	var out CreateDeviceChangeRequestResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/device-change-request", map[string]any{
		"company_id":               companyID,
		"old_device_id":            oldDeviceID,
		"device_unique_identifier": newDeviceUUID,
		"fcm_token":                newFcmToken,
		"user_agent":               newUserAgent,
		"detail":                   newDetail,
	}, bearerHeader(accessToken), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeviceChangeRequestItem represents a single device change request in a list.
type DeviceChangeRequestItem struct {
	ID                        string  `json:"id"`
	UserID                    string  `json:"user_id"`
	UserName                  string  `json:"user_name"`
	CompanyID                 string  `json:"company_id"`
	OldDeviceID               string  `json:"old_device_id"`
	NewDeviceUniqueIdentifier string  `json:"new_device_unique_identifier"`
	NewFcmToken               string  `json:"new_fcm_token"`
	NewUserAgent              string  `json:"new_user_agent"`
	NewDetail                 string  `json:"new_detail"`
	Status                    string  `json:"status"`
	RequestedByUserID         string  `json:"requested_by_user_id"`
	RequestedByName           string  `json:"requested_by_name"`
	ReviewedByUserID          *string `json:"reviewed_by_user_id,omitempty"`
	ReviewedByName            *string `json:"reviewed_by_name,omitempty"`
	ReviewedAt                *string `json:"reviewed_at,omitempty"`
	CreatedAt                 string  `json:"created_at"`
}

// DeviceChangeRequestsListResponse — GET /api/v1/device-change-requests
type DeviceChangeRequestsListResponse struct {
	Requests   []DeviceChangeRequestItem `json:"requests"`
	Page       int                       `json:"page"`
	PerPage    int                       `json:"per_page"`
	Total      int64                     `json:"total"`
	TotalPages int                       `json:"total_pages"`
}

// ListDeviceChangeRequests — GET /api/v1/device-change-requests?page=&per_page=&status=
//
// Lists all device change requests for the active company (paginated).
// status: "all" (default) | "pending" | "approved" | "rejected"
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) ListDeviceChangeRequests(accessToken, companyID string, page, perPage int, status string) (*DeviceChangeRequestsListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if status == "" {
		status = "all"
	}
	var out DeviceChangeRequestsListResponse
	_, err := c.doJSON(http.MethodGet,
		fmt.Sprintf("/api/v1/device-change-requests?page=%d&per_page=%d&status=%s", page, perPage, status),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUserPendingDeviceChangeRequestResponse — GET /api/v1/device-change-request/user-pending?user_id=
type GetUserPendingDeviceChangeRequestResponse struct {
	Request *DeviceChangeRequestItem `json:"request"`
}

// GetUserPendingDeviceChangeRequest — GET /api/v1/device-change-request/user-pending?user_id=
//
// Returns the pending device change request for a specific user in the active company.
// Returns request: null if no pending request exists.
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) GetUserPendingDeviceChangeRequest(accessToken, companyID, userID string) (*GetUserPendingDeviceChangeRequestResponse, error) {
	var out GetUserPendingDeviceChangeRequestResponse
	_, err := c.doJSON(http.MethodGet,
		fmt.Sprintf("/api/v1/device-change-request/user-pending?user_id=%s", userID),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// AcceptDeviceChangeRequest — POST /api/v1/device-change-request/:id/accept
//
// Approves a pending device change request. The old device is deleted, sessions are
// revoked, and the new device is created — all inside a single DB transaction.
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) AcceptDeviceChangeRequest(accessToken, companyID, requestID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost,
		fmt.Sprintf("/api/v1/device-change-request/%s/accept", requestID),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// RejectDeviceChangeRequest — POST /api/v1/device-change-request/:id/reject
//
// Rejects a pending device change request. No device changes are made.
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) RejectDeviceChangeRequest(accessToken, companyID, requestID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost,
		fmt.Sprintf("/api/v1/device-change-request/%s/reject", requestID),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
