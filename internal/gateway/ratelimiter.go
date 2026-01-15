package gateway

import "net/http"

type RateLimiter interface {
	Allow(route *Route, r *http.Request) bool
}
