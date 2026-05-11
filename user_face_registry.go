package authclient

import (
	"net/http"
	"time"
)

// FaceRegistryRecord represents a single face-registry entry (either company-scoped or user-only).
type FaceRegistryRecord struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CompanyID *string   `json:"company_id"`
	PhotoURL  string    `json:"photo_url"`
	Metadata  *string   `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterFaceResponse — POST /api/v1/user/face-registry
//
// Both fields may be nil depending on the registration scenario:
//   - CompanyRecord  is always populated on success.
//   - UserOnlyRecord is only populated when a new user-only duplicate was also created/updated.
type RegisterFaceResponse struct {
	CompanyRecord  *FaceRegistryRecord `json:"company_record"`
	UserOnlyRecord *FaceRegistryRecord `json:"user_only_record"`
}

// RegisterFace — POST /api/v1/user/face-registry
//
// Registers or updates a face photo for a user within a company context.
//
//   - accessToken : company-scoped access token of the actor.
//   - companyID   : company context (sent as X-Company-ID header).
//   - photoURL    : URL of the face photo to register.
//   - metadata    : optional JSON metadata from a face-recognition pipeline; pass nil to omit.
//   - targetUserID: the user whose face is being registered.
//     Pass an empty string ("") to register the actor's own face (self-registration).
//     Registering for another user requires the actor to hold an admin, subadmin, or owner role.
//
// Registration logic (mirrors server-side behaviour):
//
//   - Self-registration              → upsert both the company record and the user-only record.
//   - Admin-op, no records exist     → insert company record + insert user-only duplicate.
//   - Admin-op, user-only exists     → insert company record only; user-only record is untouched.
//   - Admin-op, company record exists → update company record only; user-only record is untouched.
func (c *FioAuthClient) RegisterFace(
	accessToken, companyID, photoURL string,
	metadata *string,
	targetUserID string,
) (*RegisterFaceResponse, error) {
	body := map[string]any{
		"photo_url": photoURL,
	}
	if targetUserID != "" {
		body["target_user_id"] = targetUserID
	}
	if metadata != nil {
		body["metadata"] = *metadata
	}

	var out RegisterFaceResponse
	_, err := c.doJSON(
		http.MethodPost,
		"/api/v1/user/face-registry",
		body,
		companyContextHeaders(accessToken, companyID),
		&out,
	)
	return &out, err
}
