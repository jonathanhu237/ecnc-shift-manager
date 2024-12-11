package application

import "github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"

func (app *Application) selfCheck() error {
	if err := app.checkBlackcoreExists(); err != nil {
		return err
	}

	return nil
}

func (app *Application) checkBlackcoreExists() error {
	exists, err := app.models.Users.CheckBlackcoreExists()
	if err != nil {
		return err
	}

	if !exists {
		user := &models.User{
			Username:     "blackcore",
			PasswordHash: "$2a$10$MLmYaAb1G7vzq.OuBtYZ1OKiUTbskRPL5jAXID3lMU2fO6dIdsdGK", // ecnc_blackcore
			Email:        "initialBlackcore@ecnc.com",
			FullName:     "初始黑心",
			Role:         "黑心",
		}

		if err := app.models.Users.InsertUser(user); err != nil {
			return err
		}

		app.logger.Warn("blackcore does not exist, create a new one with password 'ecnc_blackcore', please update it later")
	}

	return nil
}
