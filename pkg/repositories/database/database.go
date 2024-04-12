package database

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	gormDB *gorm.DB
}

// New returns a new Database that connects to a Postgres server.
func New(addr string, zeroLogger *zerolog.Logger) (*Database, error) {

	config := &gorm.Config{
		Logger: logger.New(
			zeroLogger,
			logger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	}

	gormDB, err := gorm.Open(postgres.Open(addr), config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open database connection")
	}

	return &Database{gormDB}, nil
}

func (db *Database) Transaction(ctx context.Context, txFunc func(db *Database) error) error {

	tx := db.gormDB.Begin()

	if err := txFunc(&Database{tx}); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "transaction failed")
	}

	err := tx.Commit().Error
	return errors.Wrap(err, "cannot commit transaction")
}

func columnFields(column string, fields []string) []string {
	if fields == nil {
		return nil
	}

	ret := make([]string, len(fields))
	for i, f := range fields {
		ret[i] = fmt.Sprintf("%s.%s", column, f)
	}

	return ret
}
