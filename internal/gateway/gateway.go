package gateway

import (
	"net/http"
)

type Gateway struct {
	router      Router
	ratelimiter RateLimiter
}

func NewGateway(router Router, ratelimiter RateLimiter) (*Gateway, error) {
	return &Gateway{
		router:      router,
		ratelimiter: ratelimiter,
	}, nil
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, err := g.router.Match(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if route.RateLimitPolicy != nil {
		if !g.ratelimiter.Allow(route, r) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
	}

	route.Handler.ServeHTTP(w, r)
}
