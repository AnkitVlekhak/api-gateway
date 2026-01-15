package gateway

import (
	"net/http"
)

type Gateway struct {
	router           Router
	identityResolver IdentityResolver
	ratelimiter      RateLimiter
}

func NewGateway(router Router, identityResolver IdentityResolver, ratelimiter RateLimiter) (*Gateway, error) {
	return &Gateway{
		router:           router,
		identityResolver: identityResolver,
		ratelimiter:      ratelimiter,
	}, nil
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, err := g.router.Match(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	identity, err := g.identityResolver.Resolve(r)
	if err != nil {
		http.Error(w, "Unable to resolve identity", http.StatusUnauthorized)
		return
	}

	if route.RateLimitPolicy != nil {
		if !g.ratelimiter.Allow(route, identity) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
	}

	route.Handler.ServeHTTP(w, r)
}
