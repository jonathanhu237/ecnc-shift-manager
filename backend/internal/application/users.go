package application

import (
	"fmt"
	"net/http"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/wneessen/go-mail"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required"`
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

	emailExists, err := app.models.Users.CheckEmailExists(payload.Email)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}
	if emailExists {
		app.errorResponse(w, r, errEmailExistsInCreateUser)
		return
	}

	// generate random password
	random_password := app.generateRandomPassword(12)

	// send the username and password to the e-mail
	message := mail.NewMsg()
	if err := message.From(app.config.Email.Address); err != nil {
		app.internalSeverError(w, r, err)
		return
	}
	if err := message.To(payload.Email); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	message.Subject("ECNC 假勤系统 - 您的账号信息")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf("用户名: %s, 密码: %s", payload.Username, random_password))
	if err := app.mailClient.Send(message); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// hash the password
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(random_password), bcrypt.DefaultCost)
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
