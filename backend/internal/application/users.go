package application

import (
	"net/http"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		FullName string `json:"full_name" validate:"required"`
		Role     string `json:"role" validate:"required,oneof=普通助理 资深助理 黑心"`
	}

	if err := app.readJSON(r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.validate.Struct(payload); err != nil {
		app.validateError(w, r, err)
		return
	}
	// check if username and email already exists
	userExists, err := app.models.Users.CheckUserExists(payload.Username)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}
	if userExists {
		app.errorResponse(w, r, errUsernameExistsInCreateUser)
		return
	}

	emailExists, err := app.models.Users.CheckUserExists(payload.Email)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}
	if emailExists {
		app.errorResponse(w, r, errEmailExistsInCreateUser)
		return
	}

	// hash the password
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// create the user instance
	user := &models.User{
		Username:     payload.Username,
		PasswordHash: string(passwordHashBytes),
		Email:        payload.Email,
		FullName:     payload.FullName,
		Role:         payload.Role,
	}

	// insert the user instance to the database
	if err := app.models.Users.InsertUser(user); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// return a success message
	app.successResponse(w, r, "user created successfully", user)
}
