package gateway

import (
	"net/http"
)

type APIKeyIdentityResolver struct {
	headerName string
}

func NewAPIKeyIdentityResolver(headerName string) *APIKeyIdentityResolver {
	return &APIKeyIdentityResolver{
		headerName: headerName,
	}
}

func (r *APIKeyIdentityResolver) Resolve(req *http.Request) (string, error) {
	apiKey := req.Header.Get(r.headerName)
	if apiKey == "" {
		return "", ErrAPIKeyNotFound
	}
	return apiKey, nil
}
