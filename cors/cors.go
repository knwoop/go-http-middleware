package cors

import (
	"net/http"
	"strings"

	"github.com/knwoop/go-http-middleware/adapter"
)

type Option func(option *Config)

// Config configure a CORS handler.
type Config struct {
	Origins []string
	Methods []string
	Headers []string
}

// Adapter adds CORS headers to the response.
func Adapter(options ...Option) adapter.Adapter {
	return func(next http.Handler) http.Handler {
		c := &Config{
			Origins: []string{"*"},
			Methods: []string{
				http.MethodGet,
				http.MethodPut,
				http.MethodPatch,
				http.MethodPost,
				http.MethodDelete,
				http.MethodHead,
			},
			Headers: []string{
				"Content-Type",
				"Authorization",
			},
		}
		for _, option := range options {
			option(c)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(c.Headers) == 0 {
				reqHeaders := r.Header.Get("Access-Control-Request-Headers")
				if reqHeaders != "" {
					c.Headers = append(c.Headers, reqHeaders)
				}
			}

			if r.Header.Get("Origin") != "" {
				if len(c.Origins) != 0 {
					w.Header().Set("Access-Control-Allow-Origin", strings.Join(c.Origins, ", "))
				}
				if len(c.Methods) != 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.Methods, ", "))
				}
				if len(c.Headers) != 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.Headers, ", "))
				}
			}

			if r.Method == http.MethodOptions {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
