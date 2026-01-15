package builder

import (
	"github.com/AnkitVlekhak/api-gateway/internal/config"
	"github.com/AnkitVlekhak/api-gateway/internal/gateway"
	"github.com/AnkitVlekhak/api-gateway/internal/proxy"
)

func BuildRoutes(config *config.Config) ([]*gateway.Route, error) {
	routes := []*gateway.Route{}
	for _, route := range config.Routes {
		proxy, err := proxy.NewReverseProxy(route.Backend)
		if err != nil {
			return nil, err
		}
		routes = append(routes, &gateway.Route{
			Path:    route.Path,
			Handler: proxy,
			RateLimitPolicy: &gateway.RateLimitPolicy{
				Requests: route.RateLimitPolicy.Requests,
				Window:   route.RateLimitPolicy.Window,
				KeyBy:    gateway.RateLimitKey(route.RateLimitPolicy.KeyBy),
			},
		})
	}
	return routes, nil
}
