package cache

import (
	"context"
	"encoding/json"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ErrNotFound = redis.Nil

type Cache struct {
	config *server.Config

	client *redis.Client
}

func New(ctx context.Context, config *server.Config) (*Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		Protocol: 2,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis")
	}

	cmd := client.FTCreate(ctx, "room_index",
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []any{"room:", "user_rooms:"},
		},
		&redis.FieldSchema{
			FieldName: "$.id",
			As:        "room_id",
			FieldType: redis.SearchFieldTypeTag,
		},
		&redis.FieldSchema{
			FieldName: "$.[*].id",
			As:        "user_rooms_ids",
			FieldType: redis.SearchFieldTypeTag,
		})
	if err := cmd.Err(); err != nil && err.Error() != "Index already exists" {
		return nil, errors.Wrap(err, "cannot create rooms index")
	}

	return &Cache{config, client}, nil
}

func jsonGet[T any](ctx context.Context, tx redis.Cmdable, key, path string) (res T, err error) {

	resJSON, err := tx.JSONGet(ctx, key, path).Result()
	if err != nil {
		return res, errors.Wrap(err, "cannot get session")
	}

	if len(resJSON) == 0 {
		return res, ErrNotFound
	}

	if err := json.Unmarshal([]byte(resJSON), &res); err != nil {
		return res, errors.Wrap(err, "cannot json decode response")
	}

	return res, nil
}

func ftDelete(ctx context.Context, tx redis.Cmdable, index, query string) error {

	res, err := tx.FTSearchWithArgs(ctx, index, query, &redis.FTSearchOptions{NoContent: true}).
		Result()
	if err != nil {
		return errors.Wrap(err, "cannot search index")
	}

	keys := make([]string, 0, len(res.Docs))
	for _, doc := range res.Docs {
		keys = append(keys, doc.ID)
	}

	err = tx.Del(ctx, keys...).Err()
	return errors.Wrap(err, "cannot delete keys")
}
