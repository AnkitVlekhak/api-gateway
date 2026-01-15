package gateway

import "net/http"

type NoOpRateLimiter struct {
}

func NewNoOpRateLimiter() *NoOpRateLimiter {
	return &NoOpRateLimiter{}
}

func (n *NoOpRateLimiter) Allow(route *Route, r *http.Request) bool {
	return true
}
