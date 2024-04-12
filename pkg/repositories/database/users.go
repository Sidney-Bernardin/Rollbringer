package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"

	"rollbringer/pkg/domain"
)

type gormUser struct {
	ID uuid.UUID `gorm:"column:id;default:gen_random_uuid()"`

	GoogleID string `gorm:"column:google_id"`
	Username string `gorm:"column:username"`
}

func (user *gormUser) TableName() string { return "users" }

func (user *gormUser) domain() *domain.User {
	if user != nil {
		return &domain.User{
			ID:       user.ID,
			GoogleID: user.GoogleID,
			Username: user.Username,
		}
	}
	return nil
}

func (db *Database) InsertUser(ctx context.Context, user *domain.User) error {

	userModel := gormUser{
		GoogleID: user.GoogleID,
		Username: user.Username,
	}

	result := db.gormDB.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "google_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"google_id"}),
			},
			clause.Returning{
				Columns: []clause.Column{{Name: "id"}},
			},
		).
		Create(&userModel)

	if err := result.Error; err != nil {
		return errors.Wrap(err, "cannot create user")
	}

	*user = *userModel.domain()
	return nil
}

func (db *Database) GetUser(ctx context.Context, userID uuid.UUID, userFields []string) (*domain.User, error) {

	var user gormUser
	err := db.gormDB.
		Select(userFields).
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, errors.Wrap(err, "cannot get user")
	}

	return user.domain(), nil
}

func (db *Database) GetUserByGoogleID(ctx context.Context, googleID string, userFields []string) (*domain.User, error) {

	var user gormUser
	err := db.gormDB.
		Select(userFields).
		Where("google_id = ?", googleID).
		First(&user).Error

	if err != nil {
		return nil, errors.Wrap(err, "cannot get user by google id")
	}

	return user.domain(), nil
}
