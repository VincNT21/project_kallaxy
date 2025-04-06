package models

// Response from https://openlibrary.org/search.json
type ResponseBooksSearch struct {
	NumFound int `json:"numFound"`
	Docs     []struct {
		AuthorKey        []string `json:"author_key"`
		AuthorName       []string `json:"author_name"`
		CoverEditionKey  string   `json:"cover_edition_key,omitempty"`
		CoverI           int      `json:"cover_i,omitempty"`
		FirstPublishYear int      `json:"first_publish_year,omitempty"`
		Key              string   `json:"key"`
		Title            string   `json:"title"`
		Subtitle         string   `json:"subtitle,omitempty"`
	} `json:"docs"`
}

type ResponseBookISBN struct {
	PublishDate string `json:"publish_date"`
	Authors     []struct {
		Key string `json:"key"`
	} `json:"authors"`
	Subjects    []string `json:"subjects"`
	Description struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"description"`
	Title         string   `json:"title"`
	Publishers    []string `json:"publishers"`
	Isbn13        []string `json:"isbn_13"`
	NumberOfPages int      `json:"number_of_pages"`
	FullTitle     string   `json:"full_title"`
	Covers        []int    `json:"covers"`
	Key           string   `json:"key"`
}

type ResponseBookAuthor struct {
	Bio            string   `json:"bio"`
	PersonalName   string   `json:"personal_name"`
	AlternateNames []string `json:"alternate_names"`
	Name           string   `json:"name"`
	DeathDate      string   `json:"death_date"`
	Key            string   `json:"key"`
	BirthDate      string   `json:"birth_date"`
	FullerName     string   `json:"fuller_name"`
}

// From TheMovieDB.org

type ResponseMovieSearch struct {
	Results []struct {
		BackdropPath     string `json:"backdrop_path"`
		GenreIds         []int  `json:"genre_ids"`
		ID               int    `json:"id"`
		OriginalLanguage string `json:"original_language"`
		OriginalTitle    string `json:"original_title"`
		Overview         string `json:"overview"`
		PosterPath       string `json:"poster_path"`
		ReleaseDate      string `json:"release_date"`
		Title            string `json:"title"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ResponseMovieDetails struct {
	BackdropPath string `json:"backdrop_path"`
	Genres       []struct {
		Name string `json:"name"`
	} `json:"genres"`
	ID                  int      `json:"id"`
	ImdbID              string   `json:"imdb_id"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalTitle       string   `json:"original_title"`
	Overview            string   `json:"overview"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate string `json:"release_date"`
	Runtime     int    `json:"runtime"`
	Title       string `json:"title"`
}

type ResponseTvSearch struct {
	Page    int `json:"page"`
	Results []struct {
		BackdropPath     string   `json:"backdrop_path"`
		GenreIds         []int    `json:"genre_ids"`
		ID               int      `json:"id"`
		OriginCountry    []string `json:"origin_country"`
		OriginalLanguage string   `json:"original_language"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		PosterPath       string   `json:"poster_path"`
		FirstAirDate     string   `json:"first_air_date"`
		Name             string   `json:"name"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type responseMultiSearch struct {
	Page    int `json:"page"`
	Results []struct {
		BackdropPath     string   `json:"backdrop_path"`
		ID               int      `json:"id"`
		Name             string   `json:"name,omitempty"`
		OriginalName     string   `json:"original_name,omitempty"`
		Overview         string   `json:"overview"`
		PosterPath       string   `json:"poster_path"`
		MediaType        string   `json:"media_type"`
		Adult            bool     `json:"adult"`
		OriginalLanguage string   `json:"original_language"`
		GenreIds         []int    `json:"genre_ids"`
		Popularity       float64  `json:"popularity"`
		FirstAirDate     string   `json:"first_air_date,omitempty"`
		VoteAverage      float64  `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
		OriginCountry    []string `json:"origin_country,omitempty"`
		Title            string   `json:"title,omitempty"`
		OriginalTitle    string   `json:"original_title,omitempty"`
		ReleaseDate      string   `json:"release_date,omitempty"`
		Video            bool     `json:"video,omitempty"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ResponseTvDetails struct {
	BackdropPath string `json:"backdrop_path"`
	CreatedBy    []struct {
		Name   string `json:"name"`
		Gender int    `json:"gender"`
	} `json:"created_by"`
	FirstAirDate string `json:"first_air_date"`
	Genres       []struct {
		Name string `json:"name"`
	} `json:"genres"`
	ID           int    `json:"id"`
	InProduction bool   `json:"in_production"`
	LastAirDate  string `json:"last_air_date"`
	Name         string `json:"name"`
	Networks     []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	NumberOfEpisodes    int      `json:"number_of_episodes"`
	NumberOfSeasons     int      `json:"number_of_seasons"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalName        string   `json:"original_name"`
	Overview            string   `json:"overview"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	Seasons []struct {
		AirDate      string `json:"air_date"`
		EpisodeCount int    `json:"episode_count"`
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Overview     string `json:"overview"`
		PosterPath   string `json:"poster_path"`
		SeasonNumber int    `json:"season_number"`
	} `json:"seasons"`
	Status  string `json:"status"`
	Tagline string `json:"tagline"`
}

type ResponseMovieCredits struct {
	ID   int `json:"id"`
	Cast []struct {
		Gender             int    `json:"gender"`
		ID                 int    `json:"id"`
		KnownForDepartment string `json:"known_for_department"`
		Name               string `json:"name"`
		OriginalName       string `json:"original_name"`
		ProfilePath        string `json:"profile_path"`
		Character          string `json:"character"`
		Order              int    `json:"order"`
	} `json:"cast"`
	Crew []struct {
		Gender             int    `json:"gender"`
		ID                 int    `json:"id"`
		KnownForDepartment string `json:"known_for_department"`
		Name               string `json:"name"`
		OriginalName       string `json:"original_name"`
		ProfilePath        any    `json:"profile_path"`
		Department         string `json:"department"`
		Job                string `json:"job"`
	} `json:"crew"`
}

// api.rawg.io

type ResponseVideogameSearch struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		Name      string `json:"name"`
		Playtime  int    `json:"playtime"`
		Platforms []struct {
			Platform struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"platform"`
		} `json:"platforms"`
		Stores []struct {
			Store struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"store"`
		} `json:"stores"`
		Released        string `json:"released"`
		Tba             bool   `json:"tba"`
		BackgroundImage string `json:"background_image"`
		Metacritic      int    `json:"metacritic"`
		ID              int    `json:"id"`
		Genres          []struct {
			Name string `json:"name"`
		} `json:"genres"`
	} `json:"results"`
}

type ResponseVideogameDetails struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	NameOriginal            string `json:"name_original"`
	Description             string `json:"description"`
	Metacritic              int    `json:"metacritic"`
	Released                string `json:"released"`
	Tba                     bool   `json:"tba"`
	BackgroundImage         string `json:"background_image"`
	Website                 string `json:"website"`
	Playtime                int    `json:"playtime"`
	ParentAchievementsCount int    `json:"parent_achievements_count"`
	AlternativeNames        []any  `json:"alternative_names"`
	MetacriticURL           string `json:"metacritic_url"`
	Platforms               []struct {
		Platform struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			ImageBackground string `json:"image_background"`
		} `json:"platform"`
		ReleasedAt string `json:"released_at"`
	} `json:"platforms"`
	Stores []struct {
		ID    int    `json:"id"`
		URL   string `json:"url"`
		Store struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Domain string `json:"domain"`
		} `json:"store"`
	} `json:"stores"`
	Developers []struct {
		Name            string `json:"name"`
		GamesCount      int    `json:"games_count"`
		ImageBackground string `json:"image_background"`
	} `json:"developers"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Publishers []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	EsrbRating     any    `json:"esrb_rating"`
	DescriptionRaw string `json:"description_raw"`
}

// https://boardgamegeek.com/xmlapi2/

type ResponseBoardgameSearch struct {
	Items struct {
		Item []struct {
			ID   string `json:"id"`
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished"`
		} `json:"item"`
		Total string `json:"total"`
	} `json:"items"`
}

type ResponseBoardgameSearchAlternative struct {
	Items struct {
		Item struct {
			ID   string `json:"id"`
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished"`
		} `json:"item"`
		Total string `json:"total"`
	} `json:"items"`
}

type ResponseBoardgameDetails struct {
	Items struct {
		Item struct {
			Description string `json:"description"`
			ID          string `json:"id"`
			Image       string `json:"image"`
			Link        []struct {
				ID    string `json:"id"`
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"link"`
			Maxplayers struct {
				Value string `json:"value"`
			} `json:"maxplayers"`
			Maxplaytime struct {
				Value string `json:"value"`
			} `json:"maxplaytime"`
			Minage struct {
				Value string `json:"value"`
			} `json:"minage"`
			Minplayers struct {
				Value string `json:"value"`
			} `json:"minplayers"`
			Minplaytime struct {
				Value string `json:"value"`
			} `json:"minplaytime"`
			Name []struct {
				Sortindex string `json:"sortindex"`
				Type      string `json:"type"`
				Value     string `json:"value"`
			} `json:"name"`
			Playingtime struct {
				Value string `json:"value"`
			} `json:"playingtime"`
			Thumbnail     string `json:"thumbnail"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished"`
		} `json:"item"`
	} `json:"items"`
}
