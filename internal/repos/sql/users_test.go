package sql

import (
	"log/slog"
	"testing"

	"github.com/Sidney-Bernardin/Rollbringer/internal"
	_ "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	// t.Parallel()

	var (
		ctx            = t.Context()
		container, url = createTestPostgres(t)
	)

	t.Run("get user", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx)
			require.NoError(t, err)
		})

		sql, err := New(ctx, &internal.Config{PostgresURL: url}, slog.New(slog.DiscardHandler))
		require.NoError(t, err)
		defer sql.(*repository).db.Close()

		inUser := createTestUser(t, sql.(*repository).db, "foo", nil, nil)

		outUser, err := sql.GetUser(ctx, inUser.ID)
		require.NoError(t, err)
		require.Equal(t, inUser, outUser)
	})
}
