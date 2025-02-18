package domain

import (
	"time"

	"github.com/google/uuid"
)

type Operation string

const (
	OperationError Operation = "ERROR"
)

type Event struct {
	Operation Operation `json:"operation"`
	Payload   any       `json:"payload"`
}

/////

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	GoogleID   *string     `json:"google_id" db:"google_id"`
	GoogleUser *GoogleUser `json:"google_user,omitempty" db:"google_user"`

	SpotifyID   *string      `json:"spotify_id" db:"spotify_id"`
	SpotifyUser *SpotifyUser `json:"spotify_user,omitempty" db:"spotify_user"`

	Username       string `json:"username" db:"username"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`

	Session *Session `json:"session"`
	Rooms   []Room   `json:"rooms"`
}

/////

type GoogleUser struct {
	GoogleID  string    `json:"google_id" db:"google_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	GivenName      string `json:"given_name" db:"given_name"`
	Email          string `json:"email" db:"email"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
}

/////

type SpotifyUser struct {
	SpotifyID string    `json:"spotify_id" db:"spotify_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	DisplayName    string  `json:"display_name" db:"display_name"`
	Email          string  `json:"email" db:"email"`
	ProfilePicture *string `json:"profile_picture" db:"profile_picture"`
}

/////

const OperationSession Operation = "SESSION"

type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	UserID uuid.UUID `json:"user_id" db:"user_id"`
	User   *User     `json:"user,omitempty" db:"user"`

	CSRFToken string `json:"csrf_token" db:"csrf_token"`
}

/////

const OperationGetSessionRequest Operation = "GET_SESSION_REQUEST"

type GetSessionRequest struct {
	SessionID uuid.UUID `json:"session_id"`
}

/////

const (
	OperationRoom  Operation = "ROOM"
	OperationRooms Operation = "ROOMS"
)

type Room struct {
	ID uuid.UUID `json:"id" db:"id"`

	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`
	Owner   *User     `json:"owner" db:"owner"`

	Name string `json:"name" db:"name"`
}

/////

const OperationGetRoomRequest Operation = "GET_ROOM_REQUEST"

type GetRoomRequest struct {
	RoomID uuid.UUID `json:"room_id"`
}

/////

const OperationGetRoomsRequest Operation = "GET_ROOMS_REQUEST"

type GetRoomsRequest struct {
	OwnerID uuid.UUID `json:"owner_id"`
}

/////

type Board struct {
	ID uuid.UUID `json:"id" db:"id"`

	RoomID uuid.UUID `json:"room_id" db:"room_id"`

	Name string `json:"name" db:"name"`
}
