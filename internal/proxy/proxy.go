package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) (http.Handler, error) {
	backendURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalHost := req.Host
		originalDirector(req)
		req.Host = backendURL.Host
		req.Header.Set("X-Forwarded-Host", originalHost)
		req.Header.Set("X-Origin-Gateway", "go-api-gateway")
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy, nil
}
