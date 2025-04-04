package accounts

type ViewSessionInfo struct {
	SessionID string `json:"session_id"`

	UserID   string       `json:"user_id"`
	UserInfo ViewUserInfo `json:"user_info"`

	CSRFToken string `json:"csrf_token"`
}

type ViewUserInfo struct {
	UserID string `json:"user_id"`

	GoogleID  *string `json:"google_id"`
	SpotifyID *string `json:"spotify_id"`

	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}
