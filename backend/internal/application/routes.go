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
		app.badRequest(w, r, errors.New("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.errorResponse(w, r, errMethodNotAllowed)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.loginHandler)
		r.With(app.getRequesterMiddleware).Post("/logout", app.logoutHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.getRequesterMiddleware)
		r.Route("/users", func(r chi.Router) {
			r.Use(app.authGuardMiddleware(blackcoreLevel))
			r.Post("/", app.createUserHandler)
			r.Get("/", app.getUsersHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.getUserMiddleware)
				r.Get("/", app.getUserHandler)
				r.Delete("/", app.deleteUser)
				r.Post("/update-role", app.updateUserRoleHandler)
			})
		})
		r.Route("/me", func(r chi.Router) {
			r.Get("/", app.getMyInfoHandler)
			r.Post("/update-password", app.updateMyPasswordHandler)
		})
		r.Route("/schedule-templates", func(r chi.Router) {
			r.Use(app.authGuardMiddleware(blackcoreLevel))
			r.Get("/", app.getAllScheduleTemplates)
		})
	})

	return r
}
