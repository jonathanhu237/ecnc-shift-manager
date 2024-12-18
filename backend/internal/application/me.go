package application

import (
	"errors"
	"net/http"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) getMyInfoHandler(w http.ResponseWriter, r *http.Request) {
	myInfo, ok := r.Context().Value(requesterCtxKey).(*models.User)
	if !ok {
		panic("getMyInfoHandler should be used after getRequesterMiddleware")
	}

	app.successResponse(w, r, "get my info successfully", myInfo)
}

func (app *Application) updateMyPasswordHandler(w http.ResponseWriter, r *http.Request) {
	requester, ok := r.Context().Value(requesterCtxKey).(*models.User)
	if !ok {
		panic("getMyInfoHandler should be used after getRequesterMiddleware")
	}

	var payload struct {
		OldPassword string `json:"oldPassword" validate:"required"`
		NewPassword string `json:"newPassword" validate:"required"`
	}
	if err := app.readJSON(r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.validate.Struct(payload); err != nil {
		app.validateError(w, r, err)
		return
	}

	// verify the old password
	if err := bcrypt.CompareHashAndPassword([]byte(requester.PasswordHash), []byte(payload.OldPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.errorResponse(w, r, errInvalidOldPasswordInResetPassword)
		default:
			app.internalSeverError(w, r, err)
		}
		return
	}

	// update the password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	requester.PasswordHash = string(newPasswordHash)
	if err := app.models.Users.UpdateUser(requester); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// response
	app.successResponse(w, r, "update password successfully", nil)
}
