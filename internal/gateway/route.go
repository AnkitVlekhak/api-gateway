package gateway

import "net/http"

type Route struct {
	Path    string
	Handler http.Handler

	RateLimitPolicy *RateLimitPolicy
	AuthPolicy      *AuthPolicy
	MetricPolicy    *MetricPolicy
}
