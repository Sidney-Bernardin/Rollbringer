package database

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormSession struct {
	ID uuid.UUID `gorm:"column:id;default:gen_random_uuid()"`

	UserID uuid.UUID `gorm:"column:user_id"`
	User   *gormUser

	CSRFToken string `gorm:"column:csrf_token"`
}

func (session *gormSession) TableName() string { return "sessions" }

func (session *gormSession) domain() *domain.Session {
	if session != nil {
		return &domain.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			User:      session.User.domain(),
			CSRFToken: session.CSRFToken,
		}
	}
	return nil
}

func (db *Database) UpsertSession(ctx context.Context, session *domain.Session) error {

	sessionModel := gormSession{
		UserID:    session.UserID,
		CSRFToken: session.CSRFToken,
	}

	err := db.gormDB.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}},
				UpdateAll: true,
			},
			clause.Returning{
				Columns: []clause.Column{{Name: "id"}},
			},
		).
		Create(&sessionModel).Error

	if err != nil {
		return errors.Wrap(err, "cannot upsert session")
	}

	*session = *sessionModel.domain()
	return nil
}

func (db *Database) GetSession(ctx context.Context, sessionID uuid.UUID, sessionFields, userFields []string) (*domain.Session, error) {

	var sessionModel gormSession
	q := db.gormDB.WithContext(ctx).
		Select(columnFields("sessions", sessionFields)).
		Where("id = ?", sessionID)

	if userFields != nil {
		q.Joins("User")
	}

	if err := q.First(&sessionModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.ProblemDetail{
				Type: domain.PDTypeUnauthorized,
			}
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	return sessionModel.domain(), nil
}
