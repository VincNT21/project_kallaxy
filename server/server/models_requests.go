package server

// Users
type parametersCreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Auth
type parametersLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type parametersConfirmPassword struct {
	Password string `json:"password"`
}

// Media
type parametersCreateMedium struct {
	Title     string                 `json:"title"`
	MediaType string                 `json:"media_type"`
	Creator   string                 `json:"creator"`
	PubDate   string                 `json:"pub_date"`
	ImageUrl  string                 `json:"image_url"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type parametersGetMediumByTitleAndType struct {
	Title     string `json:"title"`
	MediaType string `json:"media_type"`
}

type parametersGetMediaByType struct {
	MediaType string `json:"media_type"`
}

type parametersUpdateMedium struct {
	MediumID string                 `json:"medium_id"`
	Title    string                 `json:"title"`
	Creator  string                 `json:"creator"`
	PubDate  string                 `json:"pub_date"`
	ImageUrl string                 `json:"image_url"`
	Metadata map[string]interface{} `json:"metadata"`
}

type parametersDeleteMedium struct {
	MediumID string `json:"medium_id"`
}

// Records
type parametersCreateUserMediumRecord struct {
	MediumID  string `json:"medium_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Comments  string `json:"comments"`
}

type parametersUpdateRecord struct {
	RecordID  string `json:"record_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Comments  string `json:"comments"`
}
