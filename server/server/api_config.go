package server

import "github.com/VincNT21/kallaxy/server/internal/database"

type apiConfig struct {
	db        *database.Queries
	jwtsecret string
}

func newAPIConfig(db *database.Queries, jwtsecret string) *apiConfig {
	return &apiConfig{
		db:        db,
		jwtsecret: jwtsecret,
	}
}
