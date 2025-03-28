package server

type BookMetadata struct {
	ISBN      string   `json:"isbn"`
	Publisher string   `json:"publisher"`
	PageCount int      `json:"page_count"`
	BookType  string   `json:"book_type"`
	MainGenre string   `json:"main_genre"`
	SubGenres []string `json:"sub_genres"`
}

type MovieMetadata struct {
}

type VideogameMetadata struct {
}

type BoardgameMetadata struct {
}
