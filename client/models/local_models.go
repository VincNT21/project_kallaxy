package models

type ClientUser struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
}
