package models

type RequestPasswordReset struct {
	Message    string `json:"message"`
	ResetLink  string `json:"reset_link"`
	ResetToken string `json:"reset_token"`
	Username   string `json:"username"`
}

type RequestNewPassword struct {
	ResetToken  string `json:"token"`
	NewPassword string `json:"new_password"`
}
