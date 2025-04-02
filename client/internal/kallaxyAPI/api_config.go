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
	CreateMedia Endpoint
	GetMedia    Endpoint
	UpdateMedia Endpoint
	DeleteMedia Endpoint
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
	Search Endpoint
	ISBN   Endpoint
	Author Endpoint
}

type MoviesTvProxy struct {
	SearchMovie Endpoint
	SearchTV    Endpoint
	Search      Endpoint
	GetDetails  Endpoint
}

type VideogamesProxy struct {
	Search     Endpoint
	GetDetails Endpoint
}

type BoardgamesProxy struct {
	Search     Endpoint
	GetDetails Endpoint
}

// func InitApiConfig
