package gateway

import (
	"net/http"
)

type Gateway struct {
	router Router
}

func NewGateway(router Router) (*Gateway, error) {
	return &Gateway{
		router: router,
	}, nil
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, err := g.router.Match(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	route.Handler.ServeHTTP(w, r)
}
