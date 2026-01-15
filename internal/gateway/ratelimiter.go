package gateway

type RateLimiter interface {
	Allow(route *Route, identity string) bool
}
