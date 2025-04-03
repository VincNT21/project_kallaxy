package models

type ClientUser struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}
