package authclient

import "time"

// ──────────────────────────────────────────────
// Plain Go types returned by gRPC wrapper methods
// ──────────────────────────────────────────────

// GrpcUserBasic holds basic user data returned by the CheckUser RPC.
type GrpcUserBasic struct {
	ID              string
	Name            string
	Email           string
	PhoneCode       string
	Phone           string
	Status          string
	EmailVerifiedAt int64 // unix timestamp, 0 when not verified
	PhoneVerifiedAt int64 // unix timestamp, 0 when not verified
}

// GrpcUserCompanyRelation represents a single user↔company relationship row.
type GrpcUserCompanyRelation struct {
	UserID    string
	CompanyID string
	Role      Role
	CreatedAt time.Time // converted from unix timestamp
}

// GrpcCheckUserResult is returned by GrpcCheckUser.
type GrpcCheckUserResult struct {
	Found bool
	User  *GrpcUserBasic // nil when Found is false
}

// GrpcCheckUserCompanyRelationsResult is returned by GrpcCheckUserCompanyRelations.
type GrpcCheckUserCompanyRelationsResult struct {
	Found     bool
	Relations []GrpcUserCompanyRelation
}

// GrpcCheckUserCompanyRoleResult is returned by GrpcCheckUserCompanyRole.
type GrpcCheckUserCompanyRoleResult struct {
	Found    bool
	Relation *GrpcUserCompanyRelation // nil when Found is false
}

// GrpcGetUserAllRelationsResult is returned by GrpcGetUserAllRelations.
type GrpcGetUserAllRelationsResult struct {
	Found     bool
	Relations []GrpcUserCompanyRelation
}

// GrpcCompanyInfo holds company data returned by GetCompanyWithEndpoint.
type GrpcCompanyInfo struct {
	ID                string
	Name              string
	Email             string
	Phone             string // empty when not set
	DueDate           int64  // unix timestamp, 0 when not set
	EndpointID        string
	DeviceLoginPolicy string
	MaxDevices        int32
}

// GrpcEndpointInfo holds endpoint data returned by GetCompanyWithEndpoint.
type GrpcEndpointInfo struct {
	ID          string
	BackendMode string
	BaseURL     string
	DBDriver    string // empty when not set
	DBDSN       string // empty when not set
}

// GrpcGetCompanyWithEndpointResult is returned by GrpcGetCompanyWithEndpoint.
type GrpcGetCompanyWithEndpointResult struct {
	Found    bool
	Company  *GrpcCompanyInfo  // nil when Found is false
	Endpoint *GrpcEndpointInfo // nil when company has no endpoint
}
