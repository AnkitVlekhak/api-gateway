package gateway

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

type TrustedProxyIdentityResolver struct {
	trustedCIDRs []*net.IPNet
}

func NewTrustedProxyIdentityResolver(cidrs []string) (*TrustedProxyIdentityResolver, error) {
	resolver := &TrustedProxyIdentityResolver{}
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		resolver.trustedCIDRs = append(resolver.trustedCIDRs, ipNet)
	}
	return resolver, nil
}

func (r *TrustedProxyIdentityResolver) Resolve(req *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", errors.New("invalid remote address")
	}

	sourceIP := net.ParseIP(host)
	if sourceIP == nil {
		return "", errors.New("invalid ip address")
	}

	for _, cidr := range r.trustedCIDRs {
		if cidr.Contains(sourceIP) {
			clientIP, err := extractClientIPFromXFF(req)
			if err != nil {
				return "", err
			}
			return clientIP, nil
		}
	}

	return sourceIP.String(), nil
}

func extractClientIPFromXFF(req *http.Request) (string, error) {
	xff := req.Header.Get("X-Forwarded-For")
	if xff == "" {
		return "", errors.New("X-Forwarded-For header not found")
	}

	parts := strings.Split(xff, ",")
	clientIP := strings.TrimSpace(parts[0])

	ip := net.ParseIP(clientIP)
	if ip == nil {
		return "", errors.New("invalid X-Forwarded-For IP")
	}

	return ip.String(), nil
}
