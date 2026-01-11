package gateway

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type PrefixRouter struct {
	routes []*Route
}

func NewPrefixRouter(routes []*Route) *PrefixRouter {
	return &PrefixRouter{
		routes: routes,
	}
}

func (p *PrefixRouter) Match(r *http.Request) (*Route, error) {
	type match struct {
		route      *Route
		matchScore int
	}

	bestMatch := match{
		route:      nil,
		matchScore: -1,
	}
	for _, route := range p.routes {
		log.Println(route.Path)
		if strings.HasPrefix(r.URL.Path, route.Path) {
			match := match{
				route:      route,
				matchScore: len(route.Path),
			}
			if match.matchScore > bestMatch.matchScore {
				bestMatch = match
			}
		}
	}

	if bestMatch.route == nil {
		log.Println("no matching route found")
		return nil, errors.New("no matching route found")
	}

	return bestMatch.route, nil
}
