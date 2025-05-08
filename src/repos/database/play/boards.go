package play

import (
	"context"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"

	"rollbringer/src/domain/services/play"
	"rollbringer/src/repos/database"
)

type board struct {
	ID     uuid.UUID `db:"boards.id"`
	Name   string    `db:"boards.name"`
	Canvas []byte    `db:"boards.canvas"`

	UserIDs    []uuid.UUID                 `db:"board_users.user_ids"`
	Permisions [][]play.BoardUserPermision `db:"board_users.permisions"`
}

func (b board) Model() *play.Board {
	permisions := make(map[uuid.UUID][]play.BoardUserPermision, len(b.UserIDs))
	for i, userID := range b.UserIDs {
		permisions[userID] = b.Permisions[i]
	}

	return &play.Board{
		ID:             b.ID,
		Name:           play.BoardName(b.Name),
		Canvas:         b.Canvas,
		UserPermisions: permisions,
	}
}

func (db *playDatabase) CreateBoard(ctx context.Context, board *play.Board) error {
	err := db.Transaction(ctx, func(tx pgx.Tx) error {

		// Insert the board.
		err := database.Insert(ctx, db.Pool,
			`
				INSERT INTO play.boards (id, name, canvas)
				VALUES ($1, $2, $3)
			`,
			board.ID, board.Name, board.Canvas)

		if err != nil {
			return errors.Wrap(err, "cannot insert room")
		}

		// Copy the user-permisions into the board_users table.
		userIDs := slices.Collect(maps.Keys(board.UserPermisions))
		_, err = db.Pool.CopyFrom(ctx,
			pgx.Identifier{"play", "board_users"},
			[]string{"board_id", "user_id", "permisions"},
			pgx.CopyFromSlice(
				len(board.UserPermisions),
				func(i int) ([]any, error) {
					uID := userIDs[i]
					userPermisions := pgtype.FlatArray[string]([]string{})
					for _, p := range board.UserPermisions[uID] {
						userPermisions = append(userPermisions, string(p))
					}
					return []any{board.ID, uID, userPermisions}, nil
				}))

		return errors.Wrap(err, "cannot insert board_user")
	})

	return errors.Wrap(err, "transaction failed")
}

func (db *playDatabase) GetBoardByBoardID(ctx context.Context, boardID uuid.UUID) (*play.Board, error) {
	_, model, err := database.Get[board](ctx, db.Pool,
		`
			SELECT
				boards.id AS "boards.id",
				boards.name AS "boards.name",
				boards.canvas AS "boards.canvas",
				json_agg(board_users.user_id) AS "board_users.user_ids",
				json_agg(board_users.permisions) AS "board_users.permisions"
			FROM play.boards
			LEFT JOIN play.board_users ON boards.id = board_users.board_id
			WHERE id = $1
			GROUP BY boards.id
		`,
		boardID)

	return model, errors.Wrap(err, "cannot select board by board-ID")
}

func (db *playDatabase) GetBoardsByUserID(ctx context.Context, userID uuid.UUID) ([]*play.Board, error) {
	_, models, err := database.Gets[board](ctx, db.Pool,
		`
			SELECT
				boards.id AS "boards.id",
				boards.name AS "boards.name",
				json_agg(board_users.user_id) AS "board_users.user_ids",
				json_agg(board_users.permisions) AS "board_users.permisions"
			FROM play.boards
			LEFT JOIN play.board_users ON boards.id = board_users.board_id
			WHERE board_users.user_id = $1
			GROUP BY boards.id
		`,
		userID)

	return models, errors.Wrap(err, "cannot select boards by user-ID")
}
