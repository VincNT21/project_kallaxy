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
	c.apiClient.LastUser.RefreshToken = tokensUser.RefreshToken
	c.apiClient.LastUser.ID = tokensUser.ID
	c.apiClient.LastUser.Username = tokensUser.Username

	// Return data
	return tokensUser, nil
}
