package application

import (
	"errors"
	"log/slog"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) selfCheck() error {
	if err := app.checkConfig(); err != nil {
		return err
	}
	if err := app.checkBlackcoreExists(); err != nil {
		return err
	}

	return nil
}

func (app *Application) checkConfig() error {
	if app.config.JWT.Secret == "" {
		return errors.New("JWT secret is not set")
	}
	if app.config.Email.Address == "" {
		return errors.New("Email sender address is not set")
	}
	if app.config.Email.Password == "" {
		return errors.New("Email sender password is not set")
	}

	return nil
}

func (app *Application) checkBlackcoreExists() error {
	exists, err := app.models.Users.CheckBlackcoreExists()
	if err != nil {
		return err
	}

	random_password := app.generateRandomPassword(12)
	password_hash, err := bcrypt.GenerateFromPassword([]byte(random_password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if !exists {
		user := &models.User{
			Username:     "blackcore",
			PasswordHash: string(password_hash), // ecnc_blackcore
			Email:        "initialBlackcore@ecnc.com",
			FullName:     "初始黑心",
			Role:         "黑心",
		}

		if err := app.models.Users.InsertUser(user); err != nil {
			return err
		}

		app.logger.Warn("blackcore does not exist, create a new one", slog.String("password", random_password))
	}

	return nil
}
