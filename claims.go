package authclient

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims mirrors the JWT claims issued by the auth-service for regular users.
type CustomClaims struct {
	UserID       string    `json:"user_id,omitempty"`
	OldUserID    *int      `json:"old_user_id,omitempty"`
	CompanyID    string    `json:"company_id,omitempty"`
	OldCompanyID *int      `json:"old_company_id,omitempty"`
	Role         string    `json:"role,omitempty"`
	Platform     Platform  `json:"platform,omitempty"`
	TokenType    TokenType `json:"token_type,omitempty"`
	SID          string    `json:"sid,omitempty"`
	IsMobile     bool      `json:"is_mobile,omitempty"`
	jwt.RegisteredClaims
}

// S2SClaims mirrors the JWT claims issued for service-to-service tokens.
type S2SClaims struct {
	ServiceName string `json:"service_name"`
	TokenType   string `json:"token_type"`
	jwt.RegisteredClaims
}

// VerifyAndParseClaims verifies a user JWT and returns its CustomClaims.
// The JWKS is fetched with automatic caching. If the key id (kid) is not found,
// the cache is invalidated and retried once to handle key rotation.
func (c *FioAuthClient) VerifyAndParseClaims(tokenStr string) (*CustomClaims, error) {
	jwks, err := c.GetJWKS()
	if err != nil {
		return nil, fmt.Errorf("get jwks: %w", err)
	}

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		kid, _ := t.Header["kid"].(string)
		pub, err := findPublicKey(jwks, kid)
		if err != nil {
			// kid rotated — invalidate cache and retry once
			c.InvalidateJWKSCache()
			fresh, jerr := c.GetJWKS()
			if jerr != nil {
				return nil, err
			}
			return findPublicKey(fresh, kid)
		}
		return pub, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("token is not valid")
	}
	return claims, nil
}

// VerifyAndParseS2SClaims verifies a service-to-service JWT and returns its S2SClaims.
func (c *FioAuthClient) VerifyAndParseS2SClaims(tokenStr string) (*S2SClaims, error) {
	jwks, err := c.GetJWKS()
	if err != nil {
		return nil, fmt.Errorf("get jwks: %w", err)
	}

	claims := &S2SClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		kid, _ := t.Header["kid"].(string)
		pub, err := findPublicKey(jwks, kid)
		if err != nil {
			c.InvalidateJWKSCache()
			fresh, jerr := c.GetJWKS()
			if jerr != nil {
				return nil, err
			}
			return findPublicKey(fresh, kid)
		}
		return pub, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse s2s token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("s2s token is not valid")
	}
	return claims, nil
}
