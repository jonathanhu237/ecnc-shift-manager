package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	return r
}
