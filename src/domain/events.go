package domain

type (
	EventRoomJoined struct {
		RoomID   string                  `json:"room_id"`
		Newcomer EventRoomJoinedNewcomer `json:"newcomer"`
	}

	EventRoomJoinedNewcomer struct {
		UserID         string `json:"user_id"`
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
	}
)

type (
	EventNewBoard struct {
		BoardID string              `json:"board_id"`
		Name    string              `json:"name"`
		Users   []EventNewBoardUser `json:"users"`
	}

	EventNewBoardUser struct {
		UserID         string `json:"user_id"`
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
	}
)
