package adapter

import "net/http"

// Adapter represents an apply middlewares for http server.
type Adapter func(http.Handler) http.Handler

// Apply applies http server middlewares
func Apply(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; 0 <= i; i-- {
		h = adapters[i](h)
	}
	return h
}

// ClientAdapter represents middlewares of http RoundTrippers
type ClientAdapter func(http.RoundTripper) http.RoundTripper

// ClientApply applies http client middlewares
func ClientApply(r http.RoundTripper, adapters ...ClientAdapter) http.RoundTripper {
	for i := len(adapters) - 1; 0 <= i; i-- {
		r = adapters[i](r)
	}
	return r
}
