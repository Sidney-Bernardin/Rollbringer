package accounts

import (
	"context"

	"rollbringer/src"
)

const (
	DomainErrorTypeUsernameInvalid = "username_invalid"
	DomainErrorTypeUsernameTaken   = "username_taken"
)

type Service interface {
	UserCreate(ctx context.Context, view any, args *ArgsUserCreate) error
	UserGetByUsername(ctx context.Context, view any, username string) error
}

type service struct {
	config *src.Config

	db Database
}

func NewService(config *src.Config, db Database) Service {
	return &service{config, db}
}

func (svc *service) Run(ctx context.Context) error {
	return nil
}
