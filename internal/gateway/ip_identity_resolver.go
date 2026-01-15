package gateway

import (
	"errors"
	"net"
	"net/http"
)

type IPIdentityResolver struct {
}

func NewIPIdentityResolver() *IPIdentityResolver {
	return &IPIdentityResolver{}
}

func (i *IPIdentityResolver) Resolve(req *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", errors.New("invalid remote address")
	}
	return host, nil
}
