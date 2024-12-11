package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.loggerMiddleware)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.errorResponse(w, r, http.StatusNotFound, "route does not exist")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.errorResponse(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", app.loginHandler)
		r.Post("/users", app.createUserHandler)
	})

	return r
}
