package sql

import (
	"testing"

	"github.com/Sidney-Bernardin/Rollbringer/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func createTestPostgres(t *testing.T) (*postgres.PostgresContainer, string) {
	t.Helper()

	ctx := t.Context()

	container, err := postgres.Run(ctx, "postgres:17.5-alpine",
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithDatabase("rollbringer"),
		postgres.BasicWaitStrategies())

	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	url, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	require.NoError(t, migrate(url))
	require.NoError(t, container.Snapshot(ctx, postgres.WithSnapshotName("custom-snapshot")))

	return container, url
}

func createTestUser(t *testing.T, db *sqlx.DB, username string, googleID, spotifyID *string) *domain.User {
	t.Helper()

	user := &domain.User{
		ID:             domain.NewRandomUUID(),
		GoogleID:       googleID,
		SpotifyID:      spotifyID,
		Username:       domain.Username(username),
		ProfilePicture: "http://example.com/pic",
	}

	err := db.GetContext(t.Context(), user,
		`
			INSERT INTO rollbringer.users
			(id, google_id, spotify_id, username, profile_picture)
			VALUES ($1, $2, $3, $4, $5) RETURNING *
		`,
		user.ID, user.GoogleID, user.SpotifyID, user.Username, user.ProfilePicture)

	require.NoError(t, err)
	return user
}
