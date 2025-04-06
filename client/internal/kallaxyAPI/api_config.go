package kallaxyapi

type APIConfig struct {
	BaseURL   string
	AuthToken string
	Endpoints Endpoints
}

type Endpoints struct {
	Users         UsersEndpoints
	Media         MediaEndpoints
	Records       RecordsEndpoints
	Auth          AuthEndpoints
	PasswordReset PasswordResetEndpoints
	ExternalAPI   ExternalApiEndpoints
}

type Endpoint struct {
	Method string
	Path   string
}

type UsersEndpoints struct {
	CreateUser Endpoint
	GetUser    Endpoint
	UpdateUser Endpoint
	DeleteUser Endpoint
}

type MediaEndpoints struct {
	CreateMedia             Endpoint
	GetMediumByTitleAndType Endpoint
	GetMediaByType          Endpoint
	GetMediaWithRecords     Endpoint
	UpdateMedia             Endpoint
	DeleteMedia             Endpoint
}

type RecordsEndpoints struct {
	CreateRecord Endpoint
	GetRecord    Endpoint
	UpdateRecord Endpoint
	DeleteRecord Endpoint
}

type AuthEndpoints struct {
	Login              Endpoint
	Logout             Endpoint
	Refresh            Endpoint
	RevokeRefreshToken Endpoint
	ConfirmPassword    Endpoint
}

type PasswordResetEndpoints struct {
	RequestToken      Endpoint
	VerifyToken       Endpoint
	CreateNewPassword Endpoint
}

type ExternalApiEndpoints struct {
	Books      BooksProxy
	MoviesTV   MoviesTvProxy
	Videogames VideogamesProxy
	Boardgames BoardgamesProxy
}

type BooksProxy struct {
	Search  Endpoint
	ByISBN  Endpoint
	Author  Endpoint
	GetISBN Endpoint
}

type MoviesTvProxy struct {
	SearchMovie     Endpoint
	SearchTV        Endpoint
	Search          Endpoint
	GetDetails      Endpoint
	GetMovieCredits Endpoint
}

type VideogamesProxy struct {
	Search     Endpoint
	GetDetails Endpoint
}

type BoardgamesProxy struct {
	Search     Endpoint
	GetDetails Endpoint
}

// Initialize the API config struct, with all endpoints
func initApiConfig(baseURL string) *APIConfig {
	apiCfg := &APIConfig{
		BaseURL:   baseURL,
		AuthToken: "",
		Endpoints: Endpoints{
			Users: UsersEndpoints{
				CreateUser: Endpoint{
					Method: "POST",
					Path:   "/api/users",
				},
				GetUser: Endpoint{
					Method: "GET",
					Path:   "/api/users",
				},
				UpdateUser: Endpoint{
					Method: "PUT",
					Path:   "/api/users",
				},
				DeleteUser: Endpoint{
					Method: "DELETE",
					Path:   "/api/users",
				},
			},
			Media: MediaEndpoints{
				CreateMedia: Endpoint{
					Method: "POST",
					Path:   "/api/media",
				},
				GetMediumByTitleAndType: Endpoint{
					Method: "GET",
					Path:   "/api/media",
				},
				GetMediaByType: Endpoint{
					Method: "GET",
					Path:   "/api/media/type",
				},
				GetMediaWithRecords: Endpoint{
					Method: "GET",
					Path:   "/api/media_records",
				},
				UpdateMedia: Endpoint{
					Method: "PUT",
					Path:   "/api/media",
				},
				DeleteMedia: Endpoint{
					Method: "DELETE",
					Path:   "/api/media",
				},
			},
			Records: RecordsEndpoints{
				CreateRecord: Endpoint{
					Method: "POST",
					Path:   "/api/records",
				},
				GetRecord: Endpoint{
					Method: "GET",
					Path:   "/api/records",
				},
				UpdateRecord: Endpoint{
					Method: "PUT",
					Path:   "/api/records",
				},
				DeleteRecord: Endpoint{
					Method: "DELETE",
					Path:   "/api/records",
				},
			},
			Auth: AuthEndpoints{
				Login: Endpoint{
					Method: "POST",
					Path:   "/auth/login",
				},
				Logout: Endpoint{
					Method: "POST",
					Path:   "/auth/logout",
				},
				Refresh: Endpoint{
					Method: "POST",
					Path:   "/auth/refresh",
				},
				RevokeRefreshToken: Endpoint{
					Method: "POST",
					Path:   "/auth/revoke",
				},
				ConfirmPassword: Endpoint{
					Method: "GET",
					Path:   "/auth/login",
				},
			},
			PasswordReset: PasswordResetEndpoints{
				RequestToken: Endpoint{
					Method: "POST",
					Path:   "/auth/password_reset",
				},
				VerifyToken: Endpoint{
					Method: "GET",
					Path:   "/auth/password_reset",
				},
				CreateNewPassword: Endpoint{
					Method: "PUT",
					Path:   "/auth/password_reset",
				},
			},
			ExternalAPI: ExternalApiEndpoints{
				Books: BooksProxy{
					Search: Endpoint{
						Method: "GET",
						Path:   "/external_api/book/search",
					},
					ByISBN: Endpoint{
						Method: "GET",
						Path:   "/external_api/book/isbn",
					},
					Author: Endpoint{
						Method: "GET",
						Path:   "/external_api/book/author",
					},
					GetISBN: Endpoint{
						Method: "GET",
						Path:   "/external_api/book/search_isbn",
					},
				},
				MoviesTV: MoviesTvProxy{
					SearchMovie: Endpoint{
						Method: "GET",
						Path:   "/external_api/movie_tv/search_movie",
					},
					SearchTV: Endpoint{
						Method: "GET",
						Path:   "/external_api/movie_tv/search_tv",
					},
					Search: Endpoint{
						Method: "GET",
						Path:   "/external_api/movie_tv/search",
					},
					GetDetails: Endpoint{
						Method: "GET",
						Path:   "/external_api/movie_tv/",
					},
					GetMovieCredits: Endpoint{
						Method: "GET",
						Path:   "/external_api/movie_tv/movie_credits",
					},
				},
				Videogames: VideogamesProxy{
					Search: Endpoint{
						Method: "GET",
						Path:   "/external_api/videogame/search",
					},
					GetDetails: Endpoint{
						Method: "GET",
						Path:   "/external_api/videogame/",
					},
				},
				Boardgames: BoardgamesProxy{
					Search: Endpoint{
						Method: "GET",
						Path:   "/external_api/boardgame/search",
					},
					GetDetails: Endpoint{
						Method: "GET",
						Path:   "/external_api/boardgame",
					},
				},
			},
		},
	}

	return apiCfg
}
