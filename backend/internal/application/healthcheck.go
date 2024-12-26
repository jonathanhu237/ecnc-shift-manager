package application

import (
	"database/sql"
	"errors"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) healthCheck() error {
	if err := app.checkBlackCoreExists(); err != nil {
		return err
	}

	return nil
}

func (app *Application) checkBlackCoreExists() error {
	_, err := app.models.SelectUserByUsername(app.config.InitialAdmin.Username)
	if err == nil {
		// initial black core exists
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		app.logger.Warn("initial admin does not exist, creating a new one")
	default:
		return err // unknown error
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(app.config.InitialAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     app.config.InitialAdmin.Username,
		PasswordHash: string(password_hash),
		Email:        app.config.InitialAdmin.Email,
		FullName:     app.config.InitialAdmin.FullName,
		Role:         "黑心",
	}

	if err := app.models.InsertUser(user); err != nil {
		return err
	}

	return nil
}
