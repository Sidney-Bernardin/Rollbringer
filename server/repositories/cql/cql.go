package cql

import (
	"embed"
	"fmt"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/gocql/gocql"
	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cassandra"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
)

const version = 20250612165918

//go:embed Migrations/*.sql
var migrations embed.FS

type CQL struct {
	cluster *gocql.ClusterConfig
}

func New(config *server.Config) (*CQL, error) {

	cluster := gocql.NewCluster(config.CassandraHosts...)

	if err := migrate(config.CassandraHosts[0]); err != nil {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	return &CQL{cluster}, nil
}

func migrate(host string) error {

	migrationSrc, err := iofs.New(migrations, "Migrations")
	if err != nil {
		return errors.Wrap(err, "cannot create migration source")
	}

	m, err := gomigrate.NewWithSourceInstance("iofs", migrationSrc, fmt.Sprintf("cassandra://%s/rollbringer", host))
	if err != nil {
		return errors.Wrap(err, "cannot create migrate instance")
	}
	defer m.Close()

	if err := m.Migrate(version); err != nil && err != gomigrate.ErrNoChange {
		return errors.Wrap(err, "cannot migrate")
	}

	return nil
}
