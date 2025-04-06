package models

type ClientUser struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}

type ClientMedium struct {
	Title     string                 `json:"title"`
	MediaType string                 `json:"media_type"`
	Creator   string                 `json:"creator"`
	PubDate   string                 `json:"pubdate"`
	ImageUrl  string                 `json:"image_url"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type ShortOnlineSearchResult struct {
	Num           int
	TotalNumFound int
	Title         string
	ImageUrl      string
	PubDate       string
	ApiID         string
}
