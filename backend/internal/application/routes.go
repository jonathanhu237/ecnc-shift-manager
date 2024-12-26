package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	blackCoreLevel int32 = 3
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.handler.LoggerMiddleware)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.handler.Login)
		r.With(app.handler.GetRequesterMiddleware).Post("/logout", app.handler.Logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.handler.GetRequesterMiddleware)
		r.Route("/users", func(r chi.Router) {
			r.Use(app.handler.AuthGuardMiddleware(blackCoreLevel))
			r.Post("/", app.handler.CreateUser)
			r.Get("/", app.handler.GetAllUsers)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.handler.GetUserMiddleware)
				r.Get("/", app.handler.GetUser)
				r.Delete("/", app.handler.DeleteUser)
				r.Post("/update-role", app.handler.UpdateUserRole)
			})
		})
		r.Route("/me", func(r chi.Router) {
			r.Get("/", app.handler.GetMyInfo)
			r.Post("/update-password", app.handler.UpdateMyPassword)
		})
		r.Route("/schedule-templates", func(r chi.Router) {
			r.Use(app.handler.AuthGuardMiddleware(blackCoreLevel))
			r.Get("/", app.handler.GetAllScheduleTemplates)
			r.Post("/", app.handler.CreateScheduleTemplate)
		})
	})

	return r
}
