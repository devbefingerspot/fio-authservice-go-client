package authclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	ucpb "local/fio-auth-service-client/internal/pb/usercompany"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// ──────────────────────────────────────────────
// Internal helpers
// ──────────────────────────────────────────────

// grpcDial lazily creates and caches a single gRPC connection.
// TLS is used by default; call WithGRPCInsecure() for plaintext (dev/local).
// The connection is reused for all subsequent gRPC calls.
func (c *FioAuthClient) grpcDial() (*grpc.ClientConn, error) {
	c.grpcOnce.Do(func() {
		var creds credentials.TransportCredentials
		if c.grpcInsecure {
			creds = insecure.NewCredentials()
		} else {
			creds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
		}
		c.grpcConn, c.grpcDialErr = grpc.NewClient(
			c.grpcBaseURL,
			grpc.WithTransportCredentials(creds),
		)
	})
	return c.grpcConn, c.grpcDialErr
}

// grpcContext returns ctx enriched with the "authorization" metadata header
// when an API key is configured. The server interceptor expects this header.
func (c *FioAuthClient) grpcContext(ctx context.Context) context.Context {
	if c.grpcAPIKey == "" {
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", c.grpcAPIKey))
}

// ucClient returns an initialised UserCompanyServiceClient or an error.
func (c *FioAuthClient) ucClient() (ucpb.UserCompanyServiceClient, error) {
	conn, err := c.grpcDial()
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s: %w", c.grpcBaseURL, err)
	}
	return ucpb.NewUserCompanyServiceClient(conn), nil
}

// ──────────────────────────────────────────────
// Conversion helpers (proto → plain Go types)
// ──────────────────────────────────────────────

func pbToUserBasic(u *ucpb.UserBasic) *GrpcUserBasic {
	if u == nil {
		return nil
	}
	return &GrpcUserBasic{
		ID:        u.GetId(),
		Name:      u.GetName(),
		Email:     u.GetEmail(),
		PhoneCode: u.GetPhoneCode(),
		Phone:     u.GetPhone(),
		Status:    u.GetStatus(),
	}
}

func pbToRelation(r *ucpb.UserCompanyRelation) GrpcUserCompanyRelation {
	return GrpcUserCompanyRelation{
		UserID:    r.GetUserId(),
		CompanyID: r.GetCompanyId(),
		Role:      Role(r.GetRole()),
		CreatedAt: time.Unix(r.GetCreatedAt(), 0),
	}
}

func pbToRelations(rs []*ucpb.UserCompanyRelation) []GrpcUserCompanyRelation {
	out := make([]GrpcUserCompanyRelation, 0, len(rs))
	for _, r := range rs {
		out = append(out, pbToRelation(r))
	}
	return out
}

// ──────────────────────────────────────────────
// Public gRPC methods
// ──────────────────────────────────────────────

// GrpcCheckUser checks whether a user with the given ID exists
// and returns their basic profile data on success.
//
// Calls: /usercompany.UserCompanyService/CheckUser
func (c *FioAuthClient) GrpcCheckUser(ctx context.Context, userID string) (*GrpcCheckUserResult, error) {
	svc, err := c.ucClient()
	if err != nil {
		return nil, err
	}
	resp, err := svc.CheckUser(c.grpcContext(ctx), &ucpb.CheckUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	result := &GrpcCheckUserResult{Found: resp.GetFound()}
	if resp.GetFound() {
		result.User = pbToUserBasic(resp.GetUser())
	}
	return result, nil
}

// GrpcCheckUserCompanyRelations returns all roles that userID holds inside
// companyID. Found is false when the user has no relation to the company.
//
// Calls: /usercompany.UserCompanyService/CheckUserCompanyRelations
func (c *FioAuthClient) GrpcCheckUserCompanyRelations(ctx context.Context, userID, companyID string) (*GrpcCheckUserCompanyRelationsResult, error) {
	svc, err := c.ucClient()
	if err != nil {
		return nil, err
	}
	resp, err := svc.CheckUserCompanyRelations(
		c.grpcContext(ctx),
		&ucpb.CheckUserCompanyRelationsRequest{UserId: userID, CompanyId: companyID},
	)
	if err != nil {
		return nil, err
	}
	return &GrpcCheckUserCompanyRelationsResult{
		Found:     resp.GetFound(),
		Relations: pbToRelations(resp.GetRelations()),
	}, nil
}

// GrpcCheckUserCompanyRole checks whether userID holds a specific role inside
// companyID. role must be one of the Role constants (employee, owner, subadmin, admin).
//
// Calls: /usercompany.UserCompanyService/CheckUserCompanyRole
func (c *FioAuthClient) GrpcCheckUserCompanyRole(ctx context.Context, userID, companyID string, role Role) (*GrpcCheckUserCompanyRoleResult, error) {
	svc, err := c.ucClient()
	if err != nil {
		return nil, err
	}
	resp, err := svc.CheckUserCompanyRole(
		c.grpcContext(ctx),
		&ucpb.CheckUserCompanyRoleRequest{
			UserId:    userID,
			CompanyId: companyID,
			Role:      string(role),
		},
	)
	if err != nil {
		return nil, err
	}
	result := &GrpcCheckUserCompanyRoleResult{Found: resp.GetFound()}
	if resp.GetFound() && resp.GetRelation() != nil {
		rel := pbToRelation(resp.GetRelation())
		result.Relation = &rel
	}
	return result, nil
}

// GrpcGetUserAllRelations returns every company relation that userID holds
// across all companies.
//
// Calls: /usercompany.UserCompanyService/GetUserAllRelations
func (c *FioAuthClient) GrpcGetUserAllRelations(ctx context.Context, userID string) (*GrpcGetUserAllRelationsResult, error) {
	svc, err := c.ucClient()
	if err != nil {
		return nil, err
	}
	resp, err := svc.GetUserAllRelations(c.grpcContext(ctx), &ucpb.GetUserAllRelationsRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &GrpcGetUserAllRelationsResult{
		Found:     resp.GetFound(),
		Relations: pbToRelations(resp.GetRelations()),
	}, nil
}
