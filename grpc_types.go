package authclient

import "time"

// ──────────────────────────────────────────────
// Plain Go types returned by gRPC wrapper methods
// ──────────────────────────────────────────────

// GrpcUserBasic holds basic user data returned by the CheckUser RPC.
type GrpcUserBasic struct {
	ID        string
	Name      string
	Email     string
	PhoneCode string
	Phone     string
	Status    string
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
