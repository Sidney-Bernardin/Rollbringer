package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/database"
)

var boardViews = map[domain.BoardView]*database.View{
	domain.BoardViewAll: {
		Columns: `
			boards.id, boards.name, boards.room_id, boards.konva,
			rooms.id AS "room.id",
			rooms.name AS "room.name"`,
		Joins: `LEFT JOIN play.rooms ON rooms.id = boards.room_id`,
	},
	domain.BoardViewListItem: {
		Columns: `
			boards.id, boards.name, boards.room_id,
			rooms.id AS "room.id",
			rooms.name AS "room.name"`,
		Joins: `LEFT JOIN play.rooms ON rooms.id = boards.room_id`,
	},
}

/////

const qBoardInsert = ` 
WITH boards AS (
	INSERT INTO play.boards (name, room_id, konva)
	VALUES ($1, $2, $3)
	RETURNING *
)
SELECT %s FROM boards %s`

func (repo *playDatabaseRepository) BoardInsert(ctx context.Context, view domain.BoardView, board *domain.Board) error {
	v := boardViews[view]

	fmt.Println(string(board.Konva))
	// var k map[string]any
	// if err := json.Unmarshal(board.Konva, &k); err != nil {
	// 	return domain.Wrap(err, "cannot JSON decode konva", nil)
	// }

	err := repo.Insert(ctx, board, fmt.Sprintf(qBoardInsert, v.Columns, v.Joins),
		board.Name, board.RoomID, board.Konva)
	return domain.Wrap(err, "cannot insert board", nil)
}

/////

const qBoardsGet = ` 
SELECT %s FROM play.boards %s WHERE boards.%s = $1`

func (repo *playDatabaseRepository) BoardGet(ctx context.Context, view domain.BoardView, key string, value any) (*domain.Board, error) {
	v := boardViews[view]
	board := &domain.Board{}

	err := repo.GetOne(ctx, board, fmt.Sprintf(qBoardsGet, v.Columns, v.Joins, key), value)
	if err != nil {
		return nil, domain.Wrap(err, "cannot select board", nil)
	}

	return board, nil
}

func (repo *playDatabaseRepository) BoardsGet(ctx context.Context, view domain.BoardView, key string, value any) ([]*domain.Board, error) {
	v := boardViews[view]
	boards := []*domain.Board{}

	err := sqlx.SelectContext(ctx, repo.TX, &boards, fmt.Sprintf(qBoardsGet, v.Columns, v.Joins, key), value)
	if err != nil {
		return nil, domain.Wrap(err, "cannot select boards", nil)
	}

	return boards, nil
}
