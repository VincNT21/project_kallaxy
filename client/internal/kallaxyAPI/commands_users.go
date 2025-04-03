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

func (c *UsersClient) GetUserByID() (models.User, error) {

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Users.GetUser, nil)
	if err != nil {
		log.Printf("--ERROR-- with GetUserByID(): %v\n", err)
		return models.User{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetUserByID(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 401:
			return models.User{}, models.ErrUnauthorized
		case 500:
			return models.User{}, models.ErrServerIssue
		default:
			return models.User{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("--ERROR-- with GetUserByID(): %v\n", err)
		return models.User{}, err
	}

	// Return data
	return user, nil
}

type parametersUpdateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (c *UsersClient) UpdateUser(username, password, email string) (models.User, error) {
	params := parametersUpdateUser{
		Username: username,
		Password: password,
		Email:    email,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Users.UpdateUser, params)
	if err != nil {
		log.Printf("--ERROR-- with UpdateUser(): %v\n", err)
		return models.User{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with UpdateUser(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.User{}, models.ErrBadRequest
		case 401:
			return models.User{}, models.ErrUnauthorized
		case 409:
			return models.User{}, models.ErrConflict
		case 500:
			return models.User{}, models.ErrServerIssue
		default:
			return models.User{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("--ERROR-- with UpdateUser(): %v\n", err)
		return models.User{}, err
	}

	// Return data
	return user, nil
}

func (c *UsersClient) DeleteUser() error {

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Users.DeleteUser, nil)
	if err != nil {
		log.Printf("--ERROR-- with DeleteUser() request: %v\n", err)
		return err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with DeleteUser(). Response status code: %v\n", r.StatusCode)
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

	// Return data
	return nil
}
