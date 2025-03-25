package server

import "github.com/VincNT21/kallaxy/server/internal/database"

type apiConfig struct {
	db *database.Queries
}

func newAPIConfig(db *database.Queries) *apiConfig {
	return &apiConfig{
		db: db,
	}
}
