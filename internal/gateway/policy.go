package gateway

import "time"

type RateLimitPolicy struct {
	Requests int
	Window   time.Duration
	KeyBy    RateLimitKey
}

type RateLimitKey string

const (
	RateLimitKeyIP    RateLimitKey = "ip"
	RateLimitByAPIKey RateLimitKey = "api_key"
)

type AuthPolicy struct {
	Type  AuthType
	Scope []string
}

type AuthType string

const (
	AuthNone    AuthType = "none"
	AuthTypeJWT AuthType = "jwt"
)

type MetricPolicy struct {
	Enabled bool
	Name    string
}
