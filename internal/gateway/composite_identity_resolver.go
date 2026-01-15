package gateway

import (
	"net/http"
)

type CompositeIdentityResolver struct {
	resolvers []IdentityResolver
}

func NewCompositeIdentityResolver(resolvers ...IdentityResolver) *CompositeIdentityResolver {
	return &CompositeIdentityResolver{
		resolvers: resolvers,
	}
}

func (r *CompositeIdentityResolver) Resolve(req *http.Request) (string, error) {
	for _, resolver := range r.resolvers {
		identity, err := resolver.Resolve(req)
		if err == nil {
			return identity, nil
		}
	}
	return "", ErrIdentityNotResolved
}
