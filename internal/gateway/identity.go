package gateway

import "net/http"

type IdentityResolver interface {
	Resolve(*http.Request) (string, error)
}
