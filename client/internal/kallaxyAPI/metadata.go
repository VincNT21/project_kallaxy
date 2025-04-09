package kallaxyapi

type FieldSpec struct {
	FieldType string `json:"field_type"`
}

func InitMetadataFieldsMap() map[string][]string {
	// Metadata fields map
	metadataFieldsMap := map[string][]string{
		"book": {
			"page_count",
			"publishers",
			"isbn13",
			"isbn10",
			"subjects",
			"description",
		},
		"movie": {
			"imdb_id",
			"overview",
			"production_companies",
			"runtime",
			"genres",
			"cast",
			"original_language",
		},
		"series": {
			"overview",
			"status",
			"number_of_seasons",
			"number_of_episodes",
			"original_language",
			"number_of_episodes_per_season",
			"production_companies",
			"genres",
		},
		"videogame": {
			"description",
			"metacritic",
			"platforms",
			"publishers",
		},
		"boardgame": {
			"categories",
			"expansions",
			"implementations",
			"artists",
			"main_publishers",
			"min_players",
			"max_players",
		},
		"other": {},
	}

	return metadataFieldsMap
}

func InitMetadataFieldsSpecs() map[string]FieldSpec {
	var fieldSpecs map[string]FieldSpec

	fieldSpecs = map[string]FieldSpec{
		// Book fields
		"page_count":  {FieldType: "int"},
		"publishers":  {FieldType: "list"},
		"isbn13":      {FieldType: "string"},
		"isbn10":      {FieldType: "string"},
		"subjects":    {FieldType: "list"},
		"description": {FieldType: "string"},
		// Movie fields
		"imdb_id":              {FieldType: "string"},
		"overview":             {FieldType: "string"},
		"production_companies": {FieldType: "list"},
		"runtime":              {FieldType: "int"},
		"genres":               {FieldType: "list"},
		"cast":                 {FieldType: "list"},
		"original_language":    {FieldType: "string"},
		// Series fields
		"status":                        {FieldType: "string"},
		"number_of_seasons":             {FieldType: "int"},
		"number_of_episodes":            {FieldType: "int"},
		"number_of_episodes_per_season": {FieldType: "list"},
		// Videogame fields
		"metacritic": {FieldType: "int"},
		"platforms":  {FieldType: "list"},
		// Boardgame fields
		"categories":      {FieldType: "list"},
		"expansions":      {FieldType: "list"},
		"implementations": {FieldType: "list"},
		"artists":         {FieldType: "list"},
		"main_publishers": {FieldType: "list"},
		"min_players":     {FieldType: "int"},
		"max_players":     {FieldType: "int"},
	}

	return fieldSpecs
}
