package authclient

import (
	"fmt"
	"net/http"
)

// InviteUserAsEmployee — POST /api/v1/user-company/invite/employee
//
// Sends an invitation to userID to join the company as an employee.
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) InviteUserAsEmployee(accessToken, companyID, userID string) (*InvitationCreatedResponse, error) {
	var out InvitationCreatedResponse
	_, err := c.doJSON(http.MethodPost, "/api/v1/user-company/invite/employee", map[string]any{
		"user_id": userID,
	}, companyContextHeaders(accessToken, companyID), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// AcceptInvitation — POST /api/v1/user-company/invitation/:id/accept
//
// Accepts a pending invitation. Only the invitee can accept their own invitation.
func (c *FioAuthClient) AcceptInvitation(accessToken, invitationID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost,
		fmt.Sprintf("/api/v1/user-company/invitation/%s/accept", invitationID),
		nil,
		bearerHeader(accessToken),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// RejectInvitation — POST /api/v1/user-company/invitation/:id/reject
//
// Rejects a pending invitation. Only the invitee can reject their own invitation.
func (c *FioAuthClient) RejectInvitation(accessToken, invitationID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost,
		fmt.Sprintf("/api/v1/user-company/invitation/%s/reject", invitationID),
		nil,
		bearerHeader(accessToken),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ListInvitations — GET /api/v1/user-company/invitations
//
// Lists invitations received by the authenticated user.
// filter: "all" (default) | "pending"
func (c *FioAuthClient) ListInvitations(accessToken, filter string) (*InvitationsListResponse, error) {
	if filter == "" {
		filter = "all"
	}
	var out InvitationsListResponse
	_, err := c.doJSON(http.MethodGet,
		fmt.Sprintf("/api/v1/user-company/invitations?filter=%s", filter),
		nil,
		bearerHeader(accessToken),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCompanyInvitations — GET /api/v1/user-company/invitations/company
//
// Lists all invitations for the active company (paginated).
// status: "all" (default) | "pending" | "accepted" | "rejected"
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) GetCompanyInvitations(accessToken, companyID string, page, perPage int, status string) (*CompanyInvitationsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if status == "" {
		status = "all"
	}
	var out CompanyInvitationsResponse
	_, err := c.doJSON(http.MethodGet,
		fmt.Sprintf("/api/v1/user-company/invitations/company?page=%d&per_page=%d&status=%s", page, perPage, status),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCompanyInvitationStatusByID — GET /api/v1/user-company/invitation/:id/status
//
// Returns the status of a specific invitation within the active company.
// Requires admin, subadmin, or owner access (company context token).
func (c *FioAuthClient) GetCompanyInvitationStatusByID(accessToken, companyID, invitationID string) (*InvitationStatusResponse, error) {
	var out InvitationStatusResponse
	_, err := c.doJSON(http.MethodGet,
		fmt.Sprintf("/api/v1/user-company/invitation/%s/status", invitationID),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelInvitation — POST /api/v1/user-company/invitation/:id/cancel
//
// Cancels a pending invitation. Only admin, subadmin, or owner of the company can cancel.
func (c *FioAuthClient) CancelInvitation(accessToken, companyID, invitationID string) (*S2SMessageResponse, error) {
	var out S2SMessageResponse
	_, err := c.doJSON(http.MethodPost,
		fmt.Sprintf("/api/v1/user-company/invitation/%s/cancel", invitationID),
		nil,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
