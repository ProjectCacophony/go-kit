package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// NewRouter creates a new HTTP Server
func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return router
}

func NewHTTPServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}
