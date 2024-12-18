package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wneessen/go-mail"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		FullName string `json:"fullName" validate:"required"`
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
	if err := message.From(app.config.MailClientAddress); err != nil {
		app.internalSeverError(w, r, err)
		return
	}
	if err := message.To(payload.Email); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	mailPayload := MailPayload{
		To:      payload.Email,
		Subject: "ECNC 假勤系统 - 您的账号信息",
		Body:    fmt.Sprintf("用户名: %s, 密码: %s", payload.Username, random_password),
	}
	jsonData, err := json.Marshal(mailPayload)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.emailChan.PublishWithContext(
		ctx,
		"",
		"mail_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	); err != nil {
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

func (app *Application) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.SelectUsers()
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	app.successResponse(w, r, "get users successfully", users)
}

func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		app.internalSeverError(w, r, errors.New("getUserHandler must be used after getUserMiddleware"))
		return
	}

	app.successResponse(w, r, "get user successfully", user)
}

func (app *Application) updateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		app.internalSeverError(w, r, errors.New("getUserHandler must be used after getUserMiddleware"))
		return
	}

	var payload struct {
		Role string `json:"role" validate:"required,oneof=普通助理 资深助理 黑心"`
	}
	if err := app.readJSON(r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.validate.Struct(payload); err != nil {
		app.validateError(w, r, err)
		return
	}

	user.Role = payload.Role
	if err := app.models.Users.UpdateUser(user); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	app.successResponse(w, r, "user updated successfully", user)
}

func (app *Application) deleteUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		app.internalSeverError(w, r, errors.New("getUserHandler must be used after getUserMiddleware"))
		return
	}

	if err := app.models.Users.DeleteUser(user.ID); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	app.successResponse(w, r, "user deleted successfully", nil)
}
