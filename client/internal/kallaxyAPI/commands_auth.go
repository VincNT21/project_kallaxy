package kallaxyapi

import (
	"encoding/json"
	"log"

	"github.com/VincNT21/kallaxy/client/models"
)

type parametersLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *AuthClient) LoginUser(username, password string) (models.TokensAndUser, error) {
	params := parametersLogin{
		Username: username,
		Password: password,
	}

	// Make request
	resp, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Auth.Login, params)
	if err != nil {
		log.Printf("--ERROR-- with loginUser(): %v\n", err)
		return models.TokensAndUser{}, err
	}
	defer resp.Body.Close()

	// Decode response
	var tokensUser models.TokensAndUser
	err = json.NewDecoder(resp.Body).Decode(&tokensUser)
	if err != nil {
		log.Printf("--ERROR-- with loginUser(): %v\n", err)
		return models.TokensAndUser{}, err
	}

	// Store access token in memory
	c.apiClient.Config.AuthToken = tokensUser.AccessToken

	// Store refresh token and user data in memory
	c.apiClient.CurrentUser.RefreshToken = tokensUser.RefreshToken
	c.apiClient.CurrentUser.ID = tokensUser.ID
	c.apiClient.CurrentUser.Username = tokensUser.Username
	c.apiClient.CurrentUser.Email = tokensUser.Email

	// Return data
	log.Println("--DEBUG-- LoginUser() OK")
	return tokensUser, nil
}

func (c *AuthClient) LogoutUser() error {

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Auth.Logout, nil)
	if err != nil {
		log.Printf("--ERROR-- with LogoutUser(): %v\n", err)
		return err
	}
	defer r.Body.Close()

	// Delete CurrentUser data in memory
	c.apiClient.CurrentUser.RefreshToken = ""
	c.apiClient.CurrentUser.ID = ""
	c.apiClient.CurrentUser.Username = ""
	c.apiClient.CurrentUser.Email = ""

	// Return error
	log.Println("--DEBUG-- LogoutUser() OK")
	return nil
}

func (c *AuthClient) RefreshTokens() (models.Tokens, error) {

	// Make request
	r, err := c.apiClient.makeHttpRequestWithResfreshToken(c.apiClient.Config.Endpoints.Auth.Refresh)
	if err != nil {
		log.Printf("--ERROR-- with RefreshTokens(): %v\n", err)
		return models.Tokens{}, err
	}
	defer r.Body.Close()

	// Decode response
	var tokens models.Tokens
	err = json.NewDecoder(r.Body).Decode(&tokens)
	if err != nil {
		log.Printf("--ERROR-- with RefreshTokens(): %v\n", err)
		return models.Tokens{}, err
	}

	// Store tokens in memory
	c.apiClient.Config.AuthToken = tokens.AccessToken
	c.apiClient.CurrentUser.RefreshToken = tokens.RefreshToken

	// Return data
	log.Println("--DEBUG-- RefreshTokens() OK")
	return tokens, nil
}

func (c *AuthClient) ConfirmPassword(password string) error {
	type parameters struct {
		Password string `json:"password"`
	}
	params := parameters{
		Password: password,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Auth.ConfirmPassword, params)
	if err != nil {
		log.Printf("--ERROR-- with ConfirmPassword(): %v\n", err)
		return err
	}
	defer r.Body.Close()

	// Return ok
	log.Println("--DEBUG-- ConfirmPassword() OK")
	return nil
}

type parametersGetResetLink struct {
	Email string `json:"email"`
}

func (c *AuthClient) SendPasswordResetLink(email string) (models.RequestPasswordReset, error) {
	params := parametersGetResetLink{
		Email: email,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.PasswordReset.RequestToken, params)
	if err != nil {
		log.Printf("--ERROR-- with SendPasswordResetLink(): %v\n", err)
		return models.RequestPasswordReset{}, err
	}
	defer r.Body.Close()

	// Decode response
	var passwordLinkToken models.RequestPasswordReset
	err = json.NewDecoder(r.Body).Decode(&passwordLinkToken)
	if err != nil {
		log.Printf("--ERROR-- with SendPasswordResetLink(): %v\n", err)
		return models.RequestPasswordReset{}, err
	}

	// Return data
	log.Println("--DEBUG-- SendPasswordResetLink() OK")
	return passwordLinkToken, nil
}

func (c *AuthClient) SetNewPassword(password, token string) (models.User, error) {
	params := models.RequestNewPassword{
		ResetToken:  token,
		NewPassword: password,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.PasswordReset.CreateNewPassword, params)
	if err != nil {
		log.Printf("--ERROR-- with SetNewPassword(): %v\n", err)
		return models.User{}, err
	}
	defer r.Body.Close()

	// Decode response
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("--ERROR-- with SetNewPassword(): %v\n", err)
		return models.User{}, err
	}

	// Return data
	log.Println("--DEBUG-- SetNewPassword() OK")
	return user, nil
}
