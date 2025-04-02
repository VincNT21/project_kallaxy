package kallaxyapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VincNT21/kallaxy/client/models"
)

type parametersCreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (c *UsersClient) CreateUser(username, password, email string) (models.User, error) {
	params := parametersCreateUser{
		Username: username,
		Password: password,
		Email:    email,
	}

	// Make request
	resp, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Users.CreateUser, params)
	if err != nil {
		log.Printf("--ERROR-- with CreateUser(): %v\n", err)
		return models.User{}, err
	}
	defer resp.Body.Close()

	// Check response's status code
	if resp.StatusCode != 201 {
		log.Printf("--ERROR-- with loginUser(). Response status code: %v\n", resp.StatusCode)
		switch resp.StatusCode {
		case 400:
			return models.User{}, models.ErrBadRequest
		case 409:
			return models.User{}, models.ErrConflict
		case 500:
			return models.User{}, models.ErrServerIssue
		default:
			return models.User{}, fmt.Errorf("unknown error status code: %v", resp.StatusCode)
		}
	}

	// Decode response
	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Printf("--ERROR-- with CreateUser(): %v\n", err)
		return models.User{}, err
	}

	// Return data
	return user, nil
}
