package database

import (
	"context"

	"rollbringer/pkg/domain"
)

const qBoardInsert = ` 
WITH inserted_room AS (
	INSERT INTO play.boards (name, room_id)
	VALUES ($1, $2)
	RETURNING *
)
SELECT * FROM inserted_room`

func (repo *playDatabaseRepository) BoardInsert(ctx context.Context, board *domain.Board) error {
	err := repo.Insert(ctx, board, qBoardInsert,
		board.Name, board.RoomID)
	return domain.Wrap(err, "cannot insert board", nil)
}
