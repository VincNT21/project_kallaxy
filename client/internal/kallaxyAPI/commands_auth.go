package kallaxyapi

import (
	"encoding/json"
	"fmt"
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

	// Check response's status code
	if resp.StatusCode != 201 {
		log.Printf("--ERROR-- with loginUser(). Response status code: %v\n", resp.StatusCode)
		switch resp.StatusCode {
		case 401:
			return models.TokensAndUser{}, models.ErrUnauthorized
		case 500:
			return models.TokensAndUser{}, models.ErrServerIssue
		default:
			return models.TokensAndUser{}, fmt.Errorf("unknown error status code: %v", resp.StatusCode)
		}
	}

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

	// Check response's status code
	if r.StatusCode != 204 {
		log.Printf("--ERROR-- with LogoutUser(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 401:
			return models.ErrUnauthorized
		case 404:
			return models.ErrNotFound
		case 500:
			return models.ErrServerIssue
		default:
			return fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Delete CurrentUser data in memory
	c.apiClient.CurrentUser.RefreshToken = ""
	c.apiClient.CurrentUser.ID = ""
	c.apiClient.CurrentUser.Username = ""
	c.apiClient.CurrentUser.Email = ""

	// Return error
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

	// Check response's status code
	if r.StatusCode != 201 {
		log.Printf("--ERROR-- with RefreshTokens(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 401:
			return models.Tokens{}, models.ErrUnauthorized
		case 500:
			return models.Tokens{}, models.ErrServerIssue
		default:
			return models.Tokens{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

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

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with ConfirmPassword(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ErrBadRequest
		case 401:
			return models.ErrUnauthorized
		case 404:
			return models.ErrNotFound
		case 409:
			return models.ErrConflict
		case 500:
			return models.ErrServerIssue
		default:
			return fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Return error
	return nil
}
