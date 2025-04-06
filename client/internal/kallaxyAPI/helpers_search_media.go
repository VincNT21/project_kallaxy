package kallaxyapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *HelpersClient) SearchMediaOnExternalApiByTitle(mediaType, mediumTitle, videogamePlatform string) ([]models.ShortOnlineSearchResult, error) {
	var results []models.ShortOnlineSearchResult
	switch mediaType {
	case "book":
		response, err := c.apiClient.External.SearchForBookByTitle(mediumTitle)
		if err != nil {
			return results, err
		}
		for i, found := range response.Docs {
			results = append(results, models.ShortOnlineSearchResult{
				Num:           i + 1,
				TotalNumFound: response.NumFound,
				Title:         found.Title,
				ImageUrl:      fmt.Sprintf("https://covers.openlibrary.org/b/olid/%s-M.jpg", found.CoverEditionKey),
				PubDate:       strconv.Itoa(found.FirstPublishYear),
				ApiID:         found.Key,
			})
		}
	case "movies":
		response, err := c.apiClient.External.SearchForMovieByTitle(mediumTitle)
		if err != nil {
			return results, err
		}
		for i, found := range response.Results {
			results = append(results, models.ShortOnlineSearchResult{
				Num:           i + 1,
				TotalNumFound: response.TotalResults,
				Title:         found.Title,
				ImageUrl:      fmt.Sprintf("https://image.tmdb.org/t/p/w200%s", found.PosterPath),
				PubDate:       found.ReleaseDate,
				ApiID:         strconv.Itoa(found.ID),
			})
		}
	case "series":
		response, err := c.apiClient.External.SearchForSeriesByTitle(mediumTitle)
		if err != nil {
			return results, err
		}
		for i, found := range response.Results {
			results = append(results, models.ShortOnlineSearchResult{
				Num:           i + 1,
				TotalNumFound: response.TotalResults,
				Title:         found.Name,
				ImageUrl:      fmt.Sprintf("https://image.tmdb.org/t/p/w200%s", found.PosterPath),
				PubDate:       found.FirstAirDate,
				ApiID:         strconv.Itoa(found.ID),
			})
		}
	case "videogame":
		response, err := c.apiClient.External.SearchForVideogameOnPlatformByTitle(mediumTitle, videogamePlatform)
		if err != nil {
			return results, err
		}
		for i, found := range response.Results {
			results = append(results, models.ShortOnlineSearchResult{
				Num:           i + 1,
				TotalNumFound: response.Count,
				Title:         found.Name,
				ImageUrl:      found.BackgroundImage,
				PubDate:       found.Released,
				ApiID:         strconv.Itoa(found.ID),
			})
		}
	case "boardgame":
		response, err := c.apiClient.External.SearchForBoardgameByTitle(mediumTitle)
		if err != nil {
			return results, err
		}
		totalCount, _ := strconv.Atoi(response.Items.Total)
		for i, found := range response.Items.Item {
			results = append(results, models.ShortOnlineSearchResult{
				Num:           i + 1,
				TotalNumFound: totalCount,
				Title:         found.Name.Value,
				ImageUrl:      "for image, click GetImage",
				PubDate:       found.Yearpublished.Value,
				ApiID:         found.ID,
			})
		}
	default:
		return results, errors.New("No external API is implemented for your media type")
	}

	return results, nil
}

func (c *HelpersClient) SearchMediumDetailsOnExternalApi(mediaType, mediumID string) (models.ClientMedium, error) {
	var results models.ClientMedium
	switch mediaType {
	case "book":
		// First, need to get Book ISBN from selected work key
		bookIsbn, err := c.apiClient.Helpers.GetBookISBN(mediumID)
		if err != nil {
			return models.ClientMedium{}, err
		}
		// Make request to server proxy for details
		bookDetails, err := c.apiClient.External.GetBookDetails(bookIsbn)
		if err != nil {
			return models.ClientMedium{}, err
		}
		// Create metadata map
		metadata := make(map[string]interface{})
		metadata["page_count"] = string(bookDetails.NumberOfPages)
		metadata["publishers"] = bookDetails.Publishers
		if len(bookDetails.Isbn13) != 0 {
			metadata["isbn13"] = bookDetails.Isbn13[0]
		} else {
			metadata["isbn13"] = ""
		}
		metadata["subjects"] = bookDetails.Subjects
		metadata["description"] = bookDetails.Description.Value

		// Get author(s)
		authorsList := []string{}
		for _, author := range bookDetails.Authors {
			authorName, err := c.apiClient.Helpers.GetBookAuthor(author.Key)
			if err != nil {
				return models.ClientMedium{}, err
			}
			authorsList = append(authorsList, authorName)
		}
		authors := strings.Join(authorsList, ", ")

		// Create ClientMedium
		results = models.ClientMedium{
			Title:     bookDetails.FullTitle,
			MediaType: "book",
			Creator:   authors,
			PubDate:   bookDetails.PublishDate,
			ImageUrl:  "",
			Metadata:  metadata,
		}

	case "movie":
		// Make request to server proxy for movie details
		movieDetails, err := c.apiClient.External.GetMovieDetails(mediumID)
		if err != nil {
			return models.ClientMedium{}, err
		}
		// Make request to server proxy for movie cast details
		movieCredits, err := c.apiClient.External.GetMovieCredits(mediumID)
		if err != nil {
			return models.ClientMedium{}, err
		}

		// Create metadata map
		metadata := make(map[string]interface{})
		metadata["imdb_id"] = movieDetails.ImdbID
		metadata["overview"] = movieDetails.Overview

		productionCieList := []string{}
		for _, prodCie := range movieDetails.ProductionCompanies {
			productionCieList = append(productionCieList, prodCie.Name)
		}
		metadata["production_companies"] = productionCieList
		metadata["runtime"] = string(movieDetails.Runtime)

		genresList := []string{}
		for _, genre := range movieDetails.Genres {
			genresList = append(genresList, genre.Name)
		}
		metadata["genres"] = genresList
		metadata["cast"] = findMainCast(movieCredits)
		metadata["original_language"] = movieDetails.OriginalLanguage

		// Create ClientMedium
		results = models.ClientMedium{
			Title:     movieDetails.Title,
			MediaType: "movie",
			Creator:   findMovieDirectors(movieCredits),
			PubDate:   movieDetails.ReleaseDate,
			ImageUrl:  "",
			Metadata:  metadata,
		}

	case "series":
		// Make request to server proxy for details
		seriesDetails, err := c.apiClient.External.GetSeriesDetails(mediumID)
		if err != nil {
			return models.ClientMedium{}, err
		}

		// Create metadata map
		metadata := make(map[string]interface{})
		metadata["overview"] = seriesDetails.Overview
		metadata["status"] = seriesDetails.Status
		metadata["number_of_seasons"] = seriesDetails.NumberOfSeasons
		metadata["number_of_episodes"] = seriesDetails.NumberOfEpisodes
		metadata["original_language"] = seriesDetails.OriginalLanguage
		metadata["number_of_episodes_per_season"] = findSeasonsDetails(seriesDetails)

		productionCieList := []string{}
		for _, prodCie := range seriesDetails.ProductionCompanies {
			productionCieList = append(productionCieList, prodCie.Name)
		}
		metadata["production_companies"] = productionCieList

		genresList := []string{}
		for _, genre := range seriesDetails.Genres {
			genresList = append(genresList, genre.Name)
		}
		metadata["genres"] = genresList

		// Create ClientMedium
		results = models.ClientMedium{
			Title:     seriesDetails.Name,
			MediaType: "series",
			Creator:   findSeriesCreators(seriesDetails),
			PubDate:   seriesDetails.FirstAirDate,
			ImageUrl:  "",
			Metadata:  metadata,
		}

	case "videogame":
		// Make request to server proxy for details
		vgDetails, err := c.apiClient.External.GetVideogameDetails(mediumID)
		if err != nil {
			return models.ClientMedium{}, err
		}

		// Create metadata map
		metadata := make(map[string]interface{})
		metadata["description"] = vgDetails.DescriptionRaw
		metadata["metacritic"] = string(vgDetails.Metacritic)
		metadata["platforms"] = findVideogamePlatforms(vgDetails)

		genresList := []string{}
		for _, genre := range vgDetails.Genres {
			genresList = append(genresList, genre.Name)
		}
		metadata["genres"] = genresList

		publishersList := []string{}
		for _, publisher := range vgDetails.Publishers {
			publishersList = append(publishersList, publisher.Name)
		}
		metadata["publishers"] = publishersList

		// Create ClientMedium
		results = models.ClientMedium{
			Title:     vgDetails.Name,
			MediaType: "videogame",
			Creator:   findVideogameDevelopers(vgDetails),
			PubDate:   vgDetails.Released,
			ImageUrl:  "",
			Metadata:  metadata,
		}

	case "boardgame":
		// Make request to server proxy for details
		bgDetails, err := c.apiClient.External.GetBoardgameDetails(mediumID)
		if err != nil {
			return models.ClientMedium{}, err
		}

		// Get additional details from result "Links"
		addDetails := findBoardgameCrewAndDetails(bgDetails)

		// Create metadata map
		metadata := make(map[string]interface{})
		metadata["categories"] = addDetails["categories"]
		metadata["expansions"] = addDetails["expansions"]
		metadata["implementations"] = addDetails["implementations"]
		metadata["artists"] = addDetails["artists"]
		metadata["main_publisher"] = addDetails["main_publisher"][0]
		metadata["min_players"] = bgDetails.Items.Item.Minplayers.Value
		metadata["max_players"] = bgDetails.Items.Item.Maxplayers.Value

		// Create ClientMedium
		results = models.ClientMedium{
			Title:     bgDetails.Items.Item.Name[0].Value,
			MediaType: "boardgame",
			Creator:   strings.Join(addDetails["designers"], ", "),
			PubDate:   bgDetails.Items.Item.Yearpublished.Value,
			ImageUrl:  "",
			Metadata:  metadata,
		}

	}

	return results, nil
}

func (c *HelpersClient) SearchMediumInDB(mediaType, mediumTitle string) (models.Medium, error) {
	type parametersGetMediumByTitleAndType struct {
		Title     string `json:"title"`
		MediaType string `json:"media_type"`
	}

	// Parameters for request
	params := parametersGetMediumByTitleAndType{
		Title:     mediumTitle,
		MediaType: mediaType,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Media.GetMediumByTitleAndType, params)
	if err != nil {
		log.Printf("--ERROR-- with SearchMediumInDB(): %v\n", err)
		return models.Medium{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 201 {
		log.Printf("--ERROR-- with SearchMediumInDB(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.Medium{}, models.ErrBadRequest
		case 401:
			return models.Medium{}, models.ErrUnauthorized
		case 404:
			return models.Medium{}, models.ErrNotFound
		case 500:
			return models.Medium{}, models.ErrServerIssue
		default:
			return models.Medium{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var medium models.Medium
	err = json.NewDecoder(r.Body).Decode(&medium)
	if err != nil {
		log.Printf("--ERROR-- with SearchMediumInDB(): %v\n", err)
		return models.Medium{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchMediumInDB() OK")
	return medium, nil
}

func (c *HelpersClient) GetBoardgameImageUrl(id string) (string, error) {
	queryParameters := fmt.Sprintf("id=%s", id)

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Boardgames.GetDetails, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetBoardgameImageUrl(): %v\n", err)
		return "", err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetBoardgameImageUrl(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return "", models.ErrBadRequest
		case 401:
			return "", models.ErrUnauthorized
		case 404:
			return "", models.ErrNotFound
		case 409:
			return "", models.ErrConflict
		case 500:
			return "", models.ErrServerIssue
		default:
			return "", fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var response models.ResponseBoardgameDetails
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Printf("--ERROR-- with GetBoardgameImageUrl(): %v\n", err)
		return "", err
	}

	// Return data
	log.Println("--DEBUG-- GetBoardgameImageUrl() OK")
	return response.Items.Item.Image, nil
}

func (c *HelpersClient) GetBookISBN(worksKey string) (string, error) {

	// Create query parameters
	key := strings.TrimPrefix(worksKey, "/works/")
	queryParameters := fmt.Sprintf("key=%", key)

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Books.GetISBN, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetBookISBN(): %v\n", err)
		return "", err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetBookISBN(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return "", models.ErrBadRequest
		case 401:
			return "", models.ErrUnauthorized
		case 404:
			return "", models.ErrNotFound
		case 409:
			return "", models.ErrConflict
		case 500:
			return "", models.ErrServerIssue
		default:
			return "", fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var bookIsbn models.BookISBN
	err = json.NewDecoder(r.Body).Decode(&bookIsbn)
	if err != nil {
		log.Printf("--ERROR-- with GetBookISBN(): %v\n", err)
		return "", err
	}

	// Return data
	log.Println("--DEBUG-- GetBookISBN() OK")
	if bookIsbn.ISBN13 != "" {
		return bookIsbn.ISBN13, nil
	} else if bookIsbn.ISBN10 != "" {
		return bookIsbn.ISBN10, nil
	}

	return "", errors.New("no isbn found")
}

func (c *HelpersClient) GetBookAuthor(authorKey string) (string, error) {
	// Create query parameters
	key := strings.TrimPrefix(authorKey, "/authors/")
	queryParameters := fmt.Sprintf("author=%", key)

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Books.Author, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetBookISBN(): %v\n", err)
		return "", err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetBookISBN(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return "", models.ErrBadRequest
		case 401:
			return "", models.ErrUnauthorized
		case 404:
			return "", models.ErrNotFound
		case 409:
			return "", models.ErrConflict
		case 500:
			return "", models.ErrServerIssue
		default:
			return "", fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var bookAuthor models.ResponseBookAuthor
	err = json.NewDecoder(r.Body).Decode(&bookAuthor)
	if err != nil {
		log.Printf("--ERROR-- with GetBookISBN(): %v\n", err)
		return "", err
	}

	// Return data
	return bookAuthor.Name, nil
}

func findMovieDirectors(credits models.ResponseMovieCredits) string {
	var directorsList []string
	// Iterate over all crew to find director
	for _, crewMember := range credits.Crew {
		if crewMember.Job == "Director" {
			directorsList = append(directorsList, crewMember.Name)
		}
	}

	return strings.Join(directorsList, ", ")
}

func findSeriesCreators(seriesDetails models.ResponseTvDetails) string {
	var creatorsList []string
	// Iterate over creators
	for _, creator := range seriesDetails.CreatedBy {
		creatorsList = append(creatorsList, creator.Name)
	}

	return strings.Join(creatorsList, ", ")
}

func findMainCast(credits models.ResponseMovieCredits) []string {
	var mainCast []string
	// Get name of the three first actors/actresses of cast
	if len(credits.Cast) < 2 {
		for i := 0; i < len(credits.Cast); i++ {
			mainCast = append(mainCast, credits.Cast[i].Name)
		}
	} else {
		for _, castMember := range credits.Cast[:3] {
			mainCast = append(mainCast, castMember.Name)
		}
	}

	return mainCast
}

func findSeasonsDetails(seriesDetails models.ResponseTvDetails) []string {
	var seasonsList []string
	// Get info over each season
	for _, season := range seriesDetails.Seasons {
		seasonsList = append(seasonsList, fmt.Sprint("Season %v counts %v episodes", season.SeasonNumber, season.EpisodeCount))
	}

	return seasonsList
}

func findVideogamePlatforms(vgDetails models.ResponseVideogameDetails) []string {
	var platforms []string
	// Get info over each platform
	for _, platform := range vgDetails.Platforms {
		platforms = append(platforms, platform.Platform.Name)
	}

	return platforms
}

func findVideogameDevelopers(vgDetails models.ResponseVideogameDetails) string {
	var devs []string
	// Get info over each platform
	for _, dev := range vgDetails.Developers {
		devs = append(devs, dev.Name)
	}

	return strings.Join(devs, ", ")
}

func findBoardgameCrewAndDetails(bgDetails models.ResponseBoardgameDetails) map[string][]string {
	crew := make(map[string][]string)

	// Get info by iterating over the Links slice
	for _, link := range bgDetails.Items.Item.Link {
		switch link.Type {
		case "boardgamecategory":
			crew["categories"] = append(crew["categories"], link.Value)
		case "boardgameexpansion":
			crew["expansions"] = append(crew["expansions"], link.Value)
		case "boardgameimplementation":
			crew["implementations"] = append(crew["implementations"], link.Value)
		case "boardgamedesigner":
			crew["designers"] = append(crew["designer"], link.Value)
		case "boardgameartist":
			crew["artists"] = append(crew["artists"], link.Value)
		default:
			continue
		}
	}

	// Iterate a second time to find the first publisher only
	for _, link := range bgDetails.Items.Item.Link {
		if link.Type == "boardgamepublisher" {
			crew["original_publisher"] = append(crew["original_publisher"], link.Value)
			break
		}
	}

	return crew
}
