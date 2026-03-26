package authclient

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	gocache "github.com/patrickmn/go-cache"
)

// JWK represents a single JSON Web Key.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use,omitempty"`
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS represents a JSON Web Key Set.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

const jwksCacheKey = "jwks"

// GetJWKS fetches the JWKS from the auth server, using an in-memory cache.
func (c *FioAuthClient) GetJWKS() (*JWKS, error) {
	if cached, found := c.cache.Get(jwksCacheKey); found {
		return cached.(*JWKS), nil
	}

	var jwks JWKS
	_, err := c.doJSON(http.MethodGet, "/.well-known/jwks.json", nil, nil, &jwks)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks: %w", err)
	}

	c.cache.Set(jwksCacheKey, &jwks, gocache.DefaultExpiration)
	return &jwks, nil
}

// InvalidateJWKSCache forces a refresh on the next GetJWKS call.
func (c *FioAuthClient) InvalidateJWKSCache() {
	c.cache.Delete(jwksCacheKey)
}

// findPublicKey looks up the RSA public key for the given kid in a JWKS.
// If kid is empty, the first key is used.
func findPublicKey(jwks *JWKS, kid string) (*rsa.PublicKey, error) {
	for _, key := range jwks.Keys {
		if kid == "" || key.Kid == kid {
			return jwkToRSAPublicKey(key)
		}
	}
	return nil, fmt.Errorf("public key not found for kid: %q", kid)
}

// jwkToRSAPublicKey converts a JWK to an *rsa.PublicKey, with PKIX roundtrip
// validation to ensure the key material is well-formed.
func jwkToRSAPublicKey(key JWK) (*rsa.PublicKey, error) {
	if key.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type: %s", key.Kty)
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("decode modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	pub := &rsa.PublicKey{N: n, E: int(e.Int64())}

	// PKIX roundtrip validation
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}
	parsed, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	rsaPub, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}
