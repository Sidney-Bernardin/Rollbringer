package play

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/play/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type boardRow struct {
	ID     pgtype.UUID `db:"boards.id"`
	Name   string      `db:"boards.name"`
	Canvas []byte      `db:"boards.canvas"`

	UserIDs        []pgtype.UUID              `db:"board_users.user_ids"`
	UserPermisions [][]src.BoardUserPermision `db:"board_users.permisions"`
}

func (r *boardRow) Domain() *models.Board {
	if r == nil {
		return nil
	}

	users := make([]*src.BoardUser, len(r.UserIDs))
	for i, userID := range r.UserIDs {
		users[i] = &src.BoardUser{
			UserID:     src.UUID(userID.Bytes),
			Permisions: r.UserPermisions[i],
		}
	}

	return &models.Board{
		ID:     src.UUID(r.ID.Bytes),
		Name:   models.BoardName(r.Name),
		Canvas: r.Canvas,
		Users:  users,
	}
}

func (db *playDatabase) CreateBoard(ctx context.Context, board *models.Board) error {
	err := database.Insert(ctx, db.Tx, `
		WITH inserted_board AS (
			INSERT INTO play.boards (id, name, canvas)
			VALUES ($1, $2, $3)
		)
		INSERT INTO board_users (board_id, user_id, permisions)
		VALUES ($1, $4, $5)
	`, board.ID, board.Name, board.Canvas, board.Users[0].UserID, pq.Array(board.Users[0].Permisions))

	return errors.Wrap(err, "cannot create board")
}

func (db *playDatabase) GetBoardByBoardID(ctx context.Context, boardID src.UUID) (*models.Board, error) {
	row, err := database.Get[boardRow](ctx, db.Tx, `
        SELECT
			boards.id AS "boards.id",
			boards.name AS "boards.name",
			boards.canvas AS "boards.canvas"
        FROM play.boards
        WHERE id = $1
    `, boardID)

	return row.Domain(), errors.Wrap(err, "cannot get board by board-ID")
}

func (db *playDatabase) GetBoardsByUserID(ctx context.Context, userID src.UUID) ([]*models.Board, error) {
	rows, err := database.Gets[boardRow](ctx, db.Tx, `
		SELECT
			boards.id AS "boards.id",
			boards.name AS "boards.name",
			boards.canvas AS "boards.canvas",
			json_agg(board_users.user_id) AS "board_users.user_ids",
			json_agg(board_users.permisions) AS "board_users.permisions"
		FROM play.boards
		LEFT JOIN board_users ON boards.id = board_users.board_id
		WHERE board_users.user_id = $1
		GROUP BY boards.id
    `, userID)

	return database.Domains(rows), errors.Wrap(err, "cannot get boards by user-ID")
}
