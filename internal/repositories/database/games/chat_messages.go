package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func (db *gamesSchema) ChatMessageInsert(ctx context.Context, chatMsg *internal.ChatMessage) error {
	query := ` 
		WITH inserted_chat_message AS (
			INSERT INTO games.chat_messages (id, owner_id, game_id, message)
				VALUES ($1, $2, $3, $4)
			RETURNING *
		)
		SELECT
			inserted_chat_message.*,
			users.id AS "owner.id",
			users.username AS "owner.username",
			users.google_picture AS "owner.google_picture"
		FROM inserted_chat_message
			LEFT JOIN users.users ON users.id = inserted_chat_message.owner_id
	`

	var dbChatMsg database.ChatMessage
	err := sqlx.GetContext(ctx, db.TX, &dbChatMsg, query,
		uuid.New(), chatMsg.OwnerID, chatMsg.GameID, chatMsg.Message)

	if err != nil {
		return errors.Wrap(err, "cannot insert chat-message")
	}

	*chatMsg = *dbChatMsg.Internalized()
	return nil
}

func (db *gamesSchema) ChatMessagesGetByGame(ctx context.Context, gameID uuid.UUID) ([]*internal.ChatMessage, error) {
	query := ` 
		SELECT
			chat_messages.*,
			users.id AS "owner.id",
			users.username AS "owner.username",
			users.google_picture AS "owner.google_picture"
		FROM games.chat_messages 
			LEFT JOIN users.users ON users.id = chat_messages.owner_id
		WHERE chat_messages.game_id = $1
		ORDER BY chat_messages.created_at
	`

	var chatMsgs []*database.ChatMessage
	if err := sqlx.SelectContext(ctx, db.TX, &chatMsgs, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select chat-messages")
	}

	// Internalize each chat-message.
	ret := make([]*internal.ChatMessage, len(chatMsgs))
	for i, m := range chatMsgs {
		ret[i] = m.Internalized()
	}

	return ret, nil
}
