package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"rollbringer/pkg/domain"
)

type gormGame struct {
	ID uuid.UUID `gorm:"column:id;default:gen_random_uuid()"`

	HostID uuid.UUID `gorm:"column:host_id"`
	Host   *gormUser

	Title string `gorm:"column:title"`
}

func (game *gormGame) TableName() string { return "games" }

func (game *gormGame) domain() *domain.Game {
	if game != nil {
		return &domain.Game{
			ID:     game.ID,
			HostID: game.HostID,
			Host:   game.Host.domain(),
			Title:  game.Title,
		}
	}
	return nil
}

func (db *Database) InsertGame(ctx context.Context, game *domain.Game) error {

	gameModel := gormGame{
		HostID: game.HostID,
		Title:  game.Title,
	}

	err := db.gormDB.WithContext(ctx).Create(&gameModel).Error
	if err != nil {
		return errors.Wrap(err, "cannot create game")
	}

	*game = *gameModel.domain()
	return nil
}

func (db *Database) GetGamesCount(ctx context.Context, hostID uuid.UUID) (int, error) {

	var count int64
	err := db.gormDB.WithContext(ctx).Model(&gormGame{}).
		Where("host_id = ?", hostID).
		Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, "cannot count games")
	}

	return int(count), nil
}

func (db *Database) GetGamesByHost(ctx context.Context, hostID uuid.UUID, gameFields, hostFields []string) ([]*domain.Game, error) {

	q := db.gormDB.WithContext(ctx).
		Select(columnFields("games", gameFields)).
		Where("host_id = ?", hostID)

	if hostFields != nil {
		q = q.Joins("Host", db.gormDB.Select(hostFields))
	}

	var gameModels []gormGame
	if err := q.Find(&gameModels).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.ProblemDetail{
				Type:   domain.PDTypeGameNotFound,
				Detail: "Cannot find a game with the host-ID.",
			}
		}

		return nil, errors.Wrap(err, "cannot get games by host")
	}

	games := make([]*domain.Game, len(gameModels))
	for i, m := range gameModels {
		games[i] = m.domain()
	}

	return games, nil
}

func (db *Database) GetGame(ctx context.Context, gameID uuid.UUID, gameFields, hostFields []string) (*domain.Game, error) {

	q := db.gormDB.WithContext(ctx).
		Select(columnFields("games", gameFields)).
		Where("id = ?", gameID)

	if hostFields != nil {
		q = q.Joins("Host", db.gormDB.Select(hostFields))
	}

	var gameModel gormGame
	if err := q.First(&gameModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.ProblemDetail{
				Type:   domain.PDTypeGameNotFound,
				Detail: "Cannot find a game with the game-ID.",
			}
		}

		return nil, errors.Wrap(err, "cannot get game")
	}

	return gameModel.domain(), nil
}

func (db *Database) DeleteGame(ctx context.Context, gameID, hostID uuid.UUID) error {

	err := db.gormDB.WithContext(ctx).
		Delete(&gormGame{
			ID:     gameID,
			HostID: hostID,
		}).Error

	return errors.Wrap(err, "cannot delete game")
}
