package database

import (
	"context"
	"rollbringer/pkg/models"
)

func (database *Database) Login(ctx context.Context, user *models.User) (*models.Session, error) {
	return nil, nil
}

func (database *Database) GetSession(ctx context.Context, id string) (*models.Session, error) {
	return nil, nil
}
