package gateway

import "errors"

var (
	ErrIdentityNotResolved = errors.New("identity could not be resolved")
	ErrAPIKeyNotFound      = errors.New("API key not found")
)
