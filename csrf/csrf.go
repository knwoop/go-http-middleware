package csrf

import (
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/knwoop/go-http-middleware/adapter"
)

type Option func(option *Config)

type Config struct {
	secretKey    string
	secure       bool
	errorHandler http.Handler
}

func SetSecretKey(k string) Option {
	return func(o *Config) {
		o.secretKey = k
	}
}

func SetSecure(s bool) Option {
	return func(o *Config) {
		o.secure = s
	}
}

func ErrorHandler(h http.Handler) Option {
	return func(cs *Config) {
		cs.errorHandler = h
	}
}

func SetTokenAdapter(options ...Option) adapter.Adapter {
	return func(h http.Handler) http.Handler {
		c := &Config{
			secretKey: "32-byte-long-auth-key",
			secure:    false,
		}
		for _, option := range options {
			option(c)
		}

		csrfProtect := csrf.Protect(
			[]byte(c.secretKey),
			csrf.FieldName("csrf_token"),
			csrf.CookieName("_csrf"),
			csrf.ErrorHandler(c.errorHandler),
			csrf.Secure(c.secure),
		)
		return csrfProtect(h)
	}
}

func SetHeaderAdapter() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
			w.Header().Add("X-CSRF-Token", token)
			h.ServeHTTP(w, r)
		})
	}
}
