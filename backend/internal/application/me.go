package application

import (
	"net/http"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
)

func (app *Application) getMyInfoHandler(w http.ResponseWriter, r *http.Request) {
	myInfo, ok := r.Context().Value(requesterDetailsKey).(*models.User)
	if !ok {
		panic("getMyInfoHandler should be used after myInfoMiddleware")
	}

	app.successResponse(w, r, "get my info successfully", myInfo)
}
