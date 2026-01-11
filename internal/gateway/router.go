package gateway

import "net/http"

type Router interface {
	Match(*http.Request) (*Route, error)
}
