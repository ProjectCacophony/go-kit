package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// NewRouter creates a new HTTP Server
func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/status", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	return router
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}
