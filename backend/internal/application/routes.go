package application

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.loggerMiddleware)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w, r, errors.New("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowed(w, r)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/refresh_tokens", func(r chi.Router) {
			r.Post("/", app.createRefreshTokenHandler)
		})
	})

	return r
}
