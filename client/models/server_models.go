package models

import "encoding/json"

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokensAndUser struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Medium struct {
	ID          string          `json:"id"`
	MediaType   string          `json:"media_type"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Title       string          `json:"title"`
	Creator     string          `json:"creator"`
	ReleaseYear int32           `json:"release_year"`
	ImageUrl    string          `json:"image_url"`
	Metadata    json.RawMessage `json:"metadata"`
}

type ListMedia struct {
	Media []Medium `json:"media"`
}

type Record struct {
	ID         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	UserID     string `json:"user_id"`
	MediaID    string `json:"media_id"`
	IsFinished bool   `json:"is_finished"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Duration   int32  `json:"duration"`
}

type Records struct {
	Records []Record `json:"records"`
}

type PasswordResetRequest struct {
	Message    string `json:"message"`
	ResetLink  string `json:"reset_link"`
	ResetToken string `json:"reset_token"`
}

type ResponseVerifyResetToken struct {
	Valid bool   `json:"valid"`
	Email string `json:"email"`
}
