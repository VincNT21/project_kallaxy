package server

type ClientUser struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type ClientTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ClientTokensAndUser struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ClientMedium struct {
	ID          string                 `json:"id"`
	MediaType   string                 `json:"media_type"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Title       string                 `json:"title"`
	Creator     string                 `json:"creator"`
	ReleaseYear string                 `json:"release_year"`
	ImageUrl    string                 `json:"image_url"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type ClientListMedia struct {
	Media []ClientMedium `json:"media"`
}

type ClientRecord struct {
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

type ClientRecords struct {
	Records []ClientRecord `json:"records"`
}

// isoFormat := time.Now().UTC().Format("2006-01-02T15:04:05.999999")
