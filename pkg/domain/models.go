package domain

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type PlayPage struct {
	LoggedIn bool
	IsHost   bool

	User *User
	Game *Game
}

// =====

type UserView int

const (
	UserViewAll UserView = iota
)

var UserViews = map[string]UserView{
	"All": UserViewAll,
}

type User struct {
	ID uuid.UUID `json:"id,omitempty"`

	GoogleID *string `json:"google_id,omitempty"`
	Username string  `json:"username,omitempty"`

	PDFs        []*PDF  `json:"pdfs,omitempty"`
	HostedGames []*Game `json:"hosted_games,omitempty"`
	JoinedGames []*Game `json:"joined_games,omitempty"`
}

// =====

type SessionView int

const (
	SessionViewAll SessionView = iota
)

var SessionViews = map[string]SessionView{
	"All": SessionViewAll,
}

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID    uuid.UUID `json:"user_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

// =====

type GameView int

const (
	GameViewAll GameView = iota
	GameViewAll_HostInfo
)

var GameViews = map[string]GameView{
	"All":          GameViewAll,
	"All_HostInfo": GameViewAll_HostInfo,
}

type Game struct {
	ID uuid.UUID `json:"id,omitempty"`

	HostID uuid.UUID `json:"host_id,omitempty"`
	Host   *User     `json:"host,omitempty"`

	Name string `json:"name,omitempty"`

	Players []*User `json:"players,omitempty"`
	PDFs    []*PDF  `json:"pdfs,omitempty"`
	Rolls   []*Roll `json:"rolls,omitempty"`
}

// =====

type PDFView int

const (
	PDFViewAll PDFView = iota
	PDFViewAll_OwnerInfo_GameInfo
	PDFViewAll_OwnerInfo
	PDFViewAll_GameInfo
)

var PDFViews = map[string]PDFView{
	"All":                    PDFViewAll,
	"All_OwnerInfo_GameInfo": PDFViewAll_OwnerInfo_GameInfo,
	"All_OwnerInfo":          PDFViewAll_OwnerInfo,
	"All_GameInfo":           PDFViewAll_GameInfo,
}

type PDF struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID *uuid.UUID `json:"game_id,omitempty"`
	Game   *Game      `json:"game,omitempty"`

	Name   string              `json:"name,omitempty"`
	Schema string              `json:"schema,omitempty"`
	Fields []map[string]string `json:"fields,omitempty"`
}

// =====

type RollView int

const (
	RollViewAll RollView = iota
)

var RollViews = map[string]RollView{
	"All": RollViewAll,
}

type Roll struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID uuid.UUID `json:"game_id,omitempty"`
	Game   *Game     `json:"game,omitempty"`

	DiceNames   []int32 `json:"dice_names,omitempty"`
	DiceResults []int32 `json:"dice_results,omitempty"`
}

var validDiceNames = []int{4, 6, 8, 10, 12, 20}

func NewRoll(ctx context.Context, random *rand.Rand, ownerID, gameID uuid.UUID, diceNamesStr string) (*Roll, error) {

	roll := &Roll{
		OwnerID:     ownerID,
		GameID:      gameID,
		DiceNames:   []int32{},
		DiceResults: []int32{},
	}

	for _, dieNameStr := range strings.Split(diceNamesStr, "d")[1:] {
		dName, err := strconv.ParseInt(dieNameStr, 10, 32)
		if err != nil {
			return nil, &NormalError{
				Instance: ctx.Value(CtxKeyInstance).(string),
				Type:     NETypeInvalidDie,
				Detail:   "Dice names must resemble 32-bit integers.",
			}
		}

		roll.DiceNames = append(roll.DiceNames, int32(dName))
	}

	for _, dieName := range roll.DiceNames {
		roll.DiceResults = append(roll.DiceResults, random.Int31n(dieName)+1)
	}

	return roll, nil
}

func (r *Roll) validate(ctx context.Context) error {
	for die := range r.DiceNames {
		if !slices.Contains(validDiceNames, die) {
			return &NormalError{
				Instance: ctx.Value(CtxKeyInstance).(string),
				Type:     NETypeInvalidDie,
				Detail:   fmt.Sprintf("d%v is an invalid dice.", die),
			}
		}
	}

	return nil
}
