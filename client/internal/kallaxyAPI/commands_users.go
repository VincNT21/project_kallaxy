package kallaxyapi

import (
	"encoding/json"
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

	// Return data
	return nil
}
