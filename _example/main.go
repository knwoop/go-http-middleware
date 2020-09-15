package main

import (
	"fmt"
	"html"
	"net"
	"net/http"
	"os"

	"github.com/knwoop/go-http-middleware/cors"
	"github.com/knwoop/go-http-middleware/csrf"

	"github.com/knwoop/go-http-middleware/adapter"
)

type server struct {
	mux *http.ServeMux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) registerRoutes() {
	adapters := []adapter.Adapter{
		csrf.SetTokenAdapter(),
		csrf.SetHeaderAdapter(),
		cors.Adapter(),
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}

	s.mux.Handle("/", adapter.Apply(http.HandlerFunc(fn), adapters...))
	s.mux.Handle("/noadapter", http.HandlerFunc(fn))
}

func main() {

	s := &server{mux: http.NewServeMux()}
	s.registerRoutes()
	server := &http.Server{Handler: s}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 9090))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] net listen %s", err)
		os.Exit(1)
	}

	if err := server.Serve(lis); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] fail to serve %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
