// Package authclient provides a Go client for the Fingerspot auth-service.
//
// Usage:
//
//	client := authclient.NewFioAuthClient("http://localhost:8080", "auth-grpc.example.com:50051", "my-api-key", 30*time.Second)
//	claims, err := client.VerifyAndParseClaims(tokenString)
//
// If grpcBaseURL is an empty string, it defaults to baseURL.
// If grpcAPIKey is an empty string, gRPC calls are made without an API key.
package authclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

// FioAuthClient is the main client for interacting with the auth-service.
type FioAuthClient struct {
	baseURL      string
	grpcBaseURL  string
	grpcAPIKey   string
	grpcInsecure bool // true → plaintext; false → TLS (default)
	s2sKey       string
	httpClient   *http.Client
	cache        *gocache.Cache

	// gRPC lazy connection (initialised once on first gRPC call)
	grpcOnce    sync.Once
	grpcConn    *grpc.ClientConn
	grpcDialErr error
}

// NewFioAuthClient creates a new FioAuthClient.
//
//   - baseURL     : base URL of the auth-service HTTP API (e.g. "http://localhost:8080")
//   - grpcBaseURL : base URL for the gRPC server (e.g. "auth-grpc.example.com:50051").
//     Pass an empty string to use the same host as baseURL.
//   - grpcAPIKey  : API key sent as the "x-api-key" metadata on every gRPC call.
//     Pass an empty string to disable API key authentication.
//   - s2sKey      : pre-shared key sent as X-S2S-Authorization when issuing S2S tokens.
//   - timeout     : HTTP timeout (e.g. 30*time.Second)
//   - cacheTTL    : optional JWKS cache TTL; defaults to 5 minutes if omitted or <= 0
func NewFioAuthClient(baseURL, grpcBaseURL, grpcAPIKey, s2sKey string, timeout time.Duration, cacheTTL ...time.Duration) *FioAuthClient {
	ttl := 5 * time.Minute
	if len(cacheTTL) > 0 && cacheTTL[0] > 0 {
		ttl = cacheTTL[0]
	}
	if grpcBaseURL == "" {
		grpcBaseURL = baseURL
	}
	return &FioAuthClient{
		baseURL:     strings.TrimRight(baseURL, "/"),
		grpcBaseURL: strings.TrimRight(grpcBaseURL, "/"),
		grpcAPIKey:  grpcAPIKey,
		s2sKey:      s2sKey,
		httpClient:  &http.Client{Timeout: timeout},
		cache:       gocache.New(ttl, ttl*2),
	}
}

// HealthCheck — GET /api/v1/
func (c *FioAuthClient) HealthCheck() (*HealthCheckResponse, error) {
	var out HealthCheckResponse
	_, err := c.doJSON(http.MethodGet, "/api/v1/", nil, nil, &out)
	return &out, err
}

// doJSON performs an HTTP request and decodes the response body into out.
// On HTTP status >= 400 it returns a wrapped error containing the server's
// error message (decoded from ErrorResponse).
func (c *FioAuthClient) doJSON(method, path string, body any, headers map[string]string, out any) (int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		_ = json.Unmarshal(raw, &errResp)
		return resp.StatusCode, fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.String())
	}

	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return resp.StatusCode, fmt.Errorf("unmarshal (HTTP %d): %w | raw: %s", resp.StatusCode, err, raw)
		}
	}
	return resp.StatusCode, nil
}

// bearerHeader returns an Authorization: Bearer header map.
func bearerHeader(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

// s2sKeyHeader returns an X-S2S-Authorization header map.
func s2sKeyHeader(key string) map[string]string {
	return map[string]string{"X-S2S-Authorization": key}
}

// WithGRPCInsecure disables TLS for the gRPC connection and uses plaintext
// instead. Use this only for local development or when connecting through a
// proxy that already handles encryption (e.g. a service mesh).
//
// Must be called before the first gRPC method call.
func (c *FioAuthClient) WithGRPCInsecure() *FioAuthClient {
	c.grpcInsecure = true
	return c
}

// Close releases the cached gRPC connection. Call it when the client is no
// longer needed (e.g. in a defer after NewFioAuthClient).
func (c *FioAuthClient) Close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

// companyContextHeaders returns Authorization + X-Company-ID headers for
// endpoints that require a company context (multicompany new_web platform).
func companyContextHeaders(token, companyID string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
		"X-Company-ID":  companyID,
	}
}
