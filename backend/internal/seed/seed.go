package seed

import (
	"database/sql"
	"log/slog"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/utils"
)

type Seed struct {
	logger *slog.Logger
	config *config.Config
	models *models.Models
}

func New(logger *slog.Logger) (*Seed, *sql.DB, error) {
	seed := &Seed{
		logger: logger,
	}

	cfg, err := config.ReadConfig(logger)
	if err != nil {
		return nil, nil, err
	}
	seed.config = cfg

	db, err := utils.OpenDB(seed.config)
	if err != nil {
		return nil, nil, err
	}
	seed.models = models.New(db)

	return seed, db, nil
}

func (seed *Seed) AddRandomUsers(n int) (int, error) {
	successCnt := n

	for i := 0; i < n; i++ {
		randomUser, err := utils.GenerateRandomUser()
		if err != nil {
			return 0, err
		}

		if err := seed.models.Users.InsertUser(randomUser); err != nil {
			seed.logger.Error(
				"failed to insert user",
				slog.String("error", err.Error()),
				slog.String("username", randomUser.Username),
				slog.String("fullname", randomUser.FullName),
				slog.String("role", randomUser.Role),
			)
			n--
		}
	}

	return successCnt, nil
}
