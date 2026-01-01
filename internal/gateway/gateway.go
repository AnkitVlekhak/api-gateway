package gateway

import (
	"log"
	"net/http"

	"github.com/AnkitVlekhak/api-gateway/internal/config"
	"github.com/AnkitVlekhak/api-gateway/internal/proxy"
)

type Gateway struct {
	routes map[string]http.Handler
}

func NewGateway(config *config.Config) (*Gateway, error) {

	routes := make(map[string]http.Handler)

	for _, route := range config.Routes {
		reverseProxt, err := proxy.NewReverseProxy(route.Backend)
		if err != nil {
			return nil, err
		}
		routes[route.Path] = reverseProxt
	}

	log.Println(routes)

	return &Gateway{
		routes: routes,
	}, nil
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if handler, ok := g.routes[r.URL.Path]; ok {
		handler.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}
