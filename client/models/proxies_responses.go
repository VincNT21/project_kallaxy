package models

// Response from https://openlibrary.org/search.json
type responseBooksSearch struct {
	NumFound         int    `json:"numFound"`
	Start            int    `json:"start"`
	NumFoundExact    bool   `json:"numFoundExact"`
	NumFound0        int    `json:"num_found"`
	DocumentationURL string `json:"documentation_url"`
	Q                string `json:"q"`
	Offset           any    `json:"offset"`
	Docs             []struct {
		AuthorKey          []string `json:"author_key"`
		AuthorName         []string `json:"author_name"`
		CoverEditionKey    string   `json:"cover_edition_key,omitempty"`
		CoverI             int      `json:"cover_i,omitempty"`
		EditionCount       int      `json:"edition_count"`
		FirstPublishYear   int      `json:"first_publish_year,omitempty"`
		HasFulltext        bool     `json:"has_fulltext"`
		Ia                 []string `json:"ia,omitempty"`
		IaCollectionS      string   `json:"ia_collection_s,omitempty"`
		Key                string   `json:"key"`
		Language           []string `json:"language,omitempty"`
		LendingEditionS    string   `json:"lending_edition_s,omitempty"`
		LendingIdentifierS string   `json:"lending_identifier_s,omitempty"`
		PublicScanB        bool     `json:"public_scan_b"`
		Title              string   `json:"title"`
		Subtitle           string   `json:"subtitle,omitempty"`
	} `json:"docs"`
}

type responseBookISBN struct {
	Type struct {
		Key string `json:"key"`
	} `json:"type"`
	PublishDate    string `json:"publish_date"`
	PublishCountry string `json:"publish_country"`
	Languages      []struct {
		Key string `json:"key"`
	} `json:"languages"`
	Authors []struct {
		Key string `json:"key"`
	} `json:"authors"`
	OclcNumbers       []string `json:"oclc_numbers"`
	DeweyDecimalClass []string `json:"dewey_decimal_class"`
	WorkTitles        []string `json:"work_titles"`
	Series            []string `json:"series"`
	Notes             struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"notes"`
	Description struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"description"`
	Links []struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	} `json:"links"`
	Contributions []string `json:"contributions"`
	Subjects      []string `json:"subjects"`
	SubjectPeople []string `json:"subject_people"`
	SubjectPlaces []string `json:"subject_places"`
	Title         string   `json:"title"`
	ByStatement   string   `json:"by_statement"`
	Publishers    []string `json:"publishers"`
	PublishPlaces []string `json:"publish_places"`
	Isbn13        []string `json:"isbn_13"`
	Isbn10        []string `json:"isbn_10"`
	Pagination    string   `json:"pagination"`
	NumberOfPages int      `json:"number_of_pages"`
	Ocaid         string   `json:"ocaid"`
	SourceRecords []string `json:"source_records"`
	FullTitle     string   `json:"full_title"`
	Covers        []int    `json:"covers"`
	Works         []struct {
		Key string `json:"key"`
	} `json:"works"`
	Key         string   `json:"key"`
	LocalID     []string `json:"local_id"`
	Identifiers struct {
		Amazon []string `json:"amazon"`
	} `json:"identifiers"`
	LatestRevision int `json:"latest_revision"`
	Revision       int `json:"revision"`
	Created        struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"created"`
	LastModified struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"last_modified"`
}

type responseBookAuthor struct {
	Type struct {
		Key string `json:"key"`
	} `json:"type"`
	Links []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
		Type  struct {
			Key string `json:"key"`
		} `json:"type"`
	} `json:"links"`
	Bio            string   `json:"bio"`
	PersonalName   string   `json:"personal_name"`
	AlternateNames []string `json:"alternate_names"`
	Name           string   `json:"name"`
	DeathDate      string   `json:"death_date"`
	Key            string   `json:"key"`
	BirthDate      string   `json:"birth_date"`
	FullerName     string   `json:"fuller_name"`
	RemoteIds      struct {
		Viaf             string `json:"viaf"`
		Storygraph       string `json:"storygraph"`
		Amazon           string `json:"amazon"`
		Wikidata         string `json:"wikidata"`
		Isni             string `json:"isni"`
		Goodreads        string `json:"goodreads"`
		ProjectGutenberg string `json:"project_gutenberg"`
		Musicbrainz      string `json:"musicbrainz"`
		Bookbrainz       string `json:"bookbrainz"`
		Imdb             string `json:"imdb"`
	} `json:"remote_ids"`
	Photos         []int    `json:"photos"`
	SourceRecords  []string `json:"source_records"`
	LatestRevision int      `json:"latest_revision"`
	Revision       int      `json:"revision"`
	Created        struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"created"`
	LastModified struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"last_modified"`
}

// From TheMovieDB.org

type ResponseMovieSearch struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type responseMovieDetails struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int `json:"budget"`
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string   `json:"homepage"`
	ID                  int      `json:"id"`
	ImdbID              string   `json:"imdb_id"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalTitle       string   `json:"original_title"`
	Overview            string   `json:"overview"`
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

type ResponseTvSearch struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool     `json:"adult"`
		BackdropPath     string   `json:"backdrop_path"`
		GenreIds         []int    `json:"genre_ids"`
		ID               int      `json:"id"`
		OriginCountry    []string `json:"origin_country"`
		OriginalLanguage string   `json:"original_language"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		PosterPath       string   `json:"poster_path"`
		FirstAirDate     string   `json:"first_air_date"`
		Name             string   `json:"name"`
		VoteAverage      float64  `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
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

type responseTvDetails struct {
	Adult        bool   `json:"adult"`
	BackdropPath string `json:"backdrop_path"`
	CreatedBy    []struct {
		ID           int    `json:"id"`
		CreditID     string `json:"credit_id"`
		Name         string `json:"name"`
		OriginalName string `json:"original_name"`
		Gender       int    `json:"gender"`
		ProfilePath  string `json:"profile_path"`
	} `json:"created_by"`
	EpisodeRunTime []any  `json:"episode_run_time"`
	FirstAirDate   string `json:"first_air_date"`
	Genres         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage         string   `json:"homepage"`
	ID               int      `json:"id"`
	InProduction     bool     `json:"in_production"`
	Languages        []string `json:"languages"`
	LastAirDate      string   `json:"last_air_date"`
	LastEpisodeToAir struct {
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		EpisodeType    string  `json:"episode_type"`
		ProductionCode string  `json:"production_code"`
		Runtime        int     `json:"runtime"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int     `json:"show_id"`
		StillPath      string  `json:"still_path"`
	} `json:"last_episode_to_air"`
	Name             string `json:"name"`
	NextEpisodeToAir any    `json:"next_episode_to_air"`
	Networks         []struct {
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
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
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
		VoteAverage  int    `json:"vote_average"`
	} `json:"seasons"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Type        string  `json:"type"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// api.rawg.io

type ResponseVideogameSearch struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		Slug      string `json:"slug"`
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
		Released        string  `json:"released"`
		Tba             bool    `json:"tba"`
		BackgroundImage string  `json:"background_image"`
		Rating          float64 `json:"rating"`
		RatingTop       int     `json:"rating_top"`
		Ratings         []struct {
			ID      int     `json:"id"`
			Title   string  `json:"title"`
			Count   int     `json:"count"`
			Percent float64 `json:"percent"`
		} `json:"ratings"`
		RatingsCount     int `json:"ratings_count"`
		ReviewsTextCount int `json:"reviews_text_count"`
		Added            int `json:"added"`
		AddedByStatus    struct {
			Yet     int `json:"yet"`
			Owned   int `json:"owned"`
			Beaten  int `json:"beaten"`
			Toplay  int `json:"toplay"`
			Dropped int `json:"dropped"`
			Playing int `json:"playing"`
		} `json:"added_by_status,omitempty"`
		Metacritic       int    `json:"metacritic"`
		SuggestionsCount int    `json:"suggestions_count"`
		Updated          string `json:"updated"`
		ID               int    `json:"id"`
		Score            string `json:"score"`
		Clip             any    `json:"clip"`
		Tags             []struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Slug            string `json:"slug"`
			Language        string `json:"language"`
			GamesCount      int    `json:"games_count"`
			ImageBackground string `json:"image_background"`
		} `json:"tags"`
		EsrbRating       any    `json:"esrb_rating"`
		UserGame         any    `json:"user_game"`
		ReviewsCount     int    `json:"reviews_count"`
		SaturatedColor   string `json:"saturated_color"`
		DominantColor    string `json:"dominant_color"`
		ShortScreenshots []struct {
			ID    int    `json:"id"`
			Image string `json:"image"`
		} `json:"short_screenshots"`
		ParentPlatforms []struct {
			Platform struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"platform"`
		} `json:"parent_platforms"`
		Genres []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"genres"`
	} `json:"results"`
	UserPlatforms bool `json:"user_platforms"`
}

type responseVideogameDetails struct {
	ID                        int     `json:"id"`
	Slug                      string  `json:"slug"`
	Name                      string  `json:"name"`
	NameOriginal              string  `json:"name_original"`
	Description               string  `json:"description"`
	Metacritic                int     `json:"metacritic"`
	MetacriticPlatforms       []any   `json:"metacritic_platforms"`
	Released                  string  `json:"released"`
	Tba                       bool    `json:"tba"`
	Updated                   string  `json:"updated"`
	BackgroundImage           string  `json:"background_image"`
	BackgroundImageAdditional string  `json:"background_image_additional"`
	Website                   string  `json:"website"`
	Rating                    float64 `json:"rating"`
	RatingTop                 int     `json:"rating_top"`
	Ratings                   []struct {
		ID      int     `json:"id"`
		Title   string  `json:"title"`
		Count   int     `json:"count"`
		Percent float64 `json:"percent"`
	} `json:"ratings"`
	Reactions struct {
		Num4 int `json:"4"`
		Num6 int `json:"6"`
		Num7 int `json:"7"`
		Num8 int `json:"8"`
	} `json:"reactions"`
	Added         int `json:"added"`
	AddedByStatus struct {
		Yet     int `json:"yet"`
		Owned   int `json:"owned"`
		Beaten  int `json:"beaten"`
		Toplay  int `json:"toplay"`
		Dropped int `json:"dropped"`
		Playing int `json:"playing"`
	} `json:"added_by_status"`
	Playtime                int    `json:"playtime"`
	ScreenshotsCount        int    `json:"screenshots_count"`
	MoviesCount             int    `json:"movies_count"`
	CreatorsCount           int    `json:"creators_count"`
	AchievementsCount       int    `json:"achievements_count"`
	ParentAchievementsCount int    `json:"parent_achievements_count"`
	RedditURL               string `json:"reddit_url"`
	RedditName              string `json:"reddit_name"`
	RedditDescription       string `json:"reddit_description"`
	RedditLogo              string `json:"reddit_logo"`
	RedditCount             int    `json:"reddit_count"`
	TwitchCount             int    `json:"twitch_count"`
	YoutubeCount            int    `json:"youtube_count"`
	ReviewsTextCount        int    `json:"reviews_text_count"`
	RatingsCount            int    `json:"ratings_count"`
	SuggestionsCount        int    `json:"suggestions_count"`
	AlternativeNames        []any  `json:"alternative_names"`
	MetacriticURL           string `json:"metacritic_url"`
	ParentsCount            int    `json:"parents_count"`
	AdditionsCount          int    `json:"additions_count"`
	GameSeriesCount         int    `json:"game_series_count"`
	UserGame                any    `json:"user_game"`
	ReviewsCount            int    `json:"reviews_count"`
	SaturatedColor          string `json:"saturated_color"`
	DominantColor           string `json:"dominant_color"`
	ParentPlatforms         []struct {
		Platform struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"platform"`
	} `json:"parent_platforms"`
	Platforms []struct {
		Platform struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Slug            string `json:"slug"`
			Image           any    `json:"image"`
			YearEnd         any    `json:"year_end"`
			YearStart       any    `json:"year_start"`
			GamesCount      int    `json:"games_count"`
			ImageBackground string `json:"image_background"`
		} `json:"platform"`
		ReleasedAt   string `json:"released_at"`
		Requirements struct {
			Minimum string `json:"minimum"`
		} `json:"requirements"`
	} `json:"platforms"`
	Stores []struct {
		ID    int    `json:"id"`
		URL   string `json:"url"`
		Store struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Slug            string `json:"slug"`
			Domain          string `json:"domain"`
			GamesCount      int    `json:"games_count"`
			ImageBackground string `json:"image_background"`
		} `json:"store"`
	} `json:"stores"`
	Developers []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Slug            string `json:"slug"`
		GamesCount      int    `json:"games_count"`
		ImageBackground string `json:"image_background"`
	} `json:"developers"`
	Genres []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Slug            string `json:"slug"`
		GamesCount      int    `json:"games_count"`
		ImageBackground string `json:"image_background"`
	} `json:"genres"`
	Tags []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Slug            string `json:"slug"`
		Language        string `json:"language"`
		GamesCount      int    `json:"games_count"`
		ImageBackground string `json:"image_background"`
	} `json:"tags"`
	Publishers []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Slug            string `json:"slug"`
		GamesCount      int    `json:"games_count"`
		ImageBackground string `json:"image_background"`
	} `json:"publishers"`
	EsrbRating     any    `json:"esrb_rating"`
	Clip           any    `json:"clip"`
	DescriptionRaw string `json:"description_raw"`
}

// https://boardgamegeek.com/xmlapi2/

type ResponseBoardgameSearch struct {
	Items struct {
		Item []struct {
			ID   string `json:"id"`
			Name struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"name"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished,omitempty"`
		} `json:"item"`
		Termsofuse string `json:"termsofuse"`
		Total      string `json:"total"`
	} `json:"items"`
}

type ResponseBoardgameSearchAlternative struct {
	Items struct {
		Item struct {
			ID   string `json:"id"`
			Name struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"name"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished"`
		} `json:"item"`
		Termsofuse string `json:"termsofuse"`
		Total      string `json:"total"`
	} `json:"items"`
}

type ResponseBoardgameDetails struct {
	Items struct {
		Item struct {
			Description string `json:"description"`
			ID          string `json:"id"`
			Image       string `json:"image"`
			Maxplayers  struct {
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
			Playingtime struct {
				Value string `json:"value"`
			} `json:"playingtime"`
			Thumbnail     string `json:"thumbnail"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished"`
		} `json:"item"`
		Termsofuse string `json:"termsofuse"`
	} `json:"items"`
}
