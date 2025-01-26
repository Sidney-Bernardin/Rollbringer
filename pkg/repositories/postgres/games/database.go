package database

import (
	"io/fs"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/games"
	"rollbringer/pkg/repositories/postgres"
)

type GamesDatabaseRepository struct {
	*postgres.DatabaseRepository
}

func NewGamesDatabaseRepository(config *domain.Config, migrations fs.FS) (service.GamesDatabaseRepository, error) {

	// Create base database repository.
	dbRepo, err := postgres.NewDatabaseRepository(config, migrations)
	if err != nil {
		return nil, domain.Wrap(err, "cannot create base database repository", nil)
	}

	return &GamesDatabaseRepository{dbRepo}, nil
}
