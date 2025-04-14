package server

import "github.com/VincNT21/kallaxy/server/internal/database"

type apiConfig struct {
	db            *database.Queries
	jwtsecret     string
	openlibraryUA string
	moviedbKey    string
	rawgKey       string
	serverVersion string
}

func newAPIConfig(db *database.Queries, jwtsecret, openLibraryUA, moviedbAPIKey, rawgKey, serverVersion string) *apiConfig {
	return &apiConfig{
		db:            db,
		jwtsecret:     jwtsecret,
		openlibraryUA: openLibraryUA,
		moviedbKey:    moviedbAPIKey,
		rawgKey:       rawgKey,
		serverVersion: serverVersion,
	}
}
