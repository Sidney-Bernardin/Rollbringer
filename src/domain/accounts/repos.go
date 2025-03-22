package accounts

import (
	"context"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface {
		UserCreate(ctx context.Context, view any, cmd *CmdUserCreate) error
	}

	DatabaseQueries interface {
		UserGetByUsername(ctx context.Context, view any, username Username) error
	}
)
