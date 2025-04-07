package kallaxyapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *ExternalAPIClient) SearchForBookByTitle(bookTitle string) (models.ResponseBooksSearch, error) {
	queryParameters := fmt.Sprintf("title=%s", url.QueryEscape(bookTitle))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Books.Search, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBookByTitle(): %v\n", err)
		return models.ResponseBooksSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForBookByTitle(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseBooksSearch{}, models.ErrBadRequest
		case 401:
			return models.ResponseBooksSearch{}, models.ErrUnauthorized
		case 404:
			return models.ResponseBooksSearch{}, models.ErrNotFound
		case 409:
			return models.ResponseBooksSearch{}, models.ErrConflict
		case 500:
			return models.ResponseBooksSearch{}, models.ErrServerIssue
		default:
			return models.ResponseBooksSearch{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var books models.ResponseBooksSearch
	err = json.NewDecoder(r.Body).Decode(&books)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBookByTitle(): %v\n", err)
		return models.ResponseBooksSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForBookByTitle() OK")
	return books, nil
}

func (c *ExternalAPIClient) GetBookDetails(isbn string) (models.ResponseBookISBN, error) {
	queryParameters := fmt.Sprintf("isbn=%s", url.QueryEscape(isbn))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Books.ByISBN, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetBookDetails(): %v\n", err)
		return models.ResponseBookISBN{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetBookDetails(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseBookISBN{}, models.ErrBadRequest
		case 401:
			return models.ResponseBookISBN{}, models.ErrUnauthorized
		case 404:
			return models.ResponseBookISBN{}, models.ErrNotFound
		case 409:
			return models.ResponseBookISBN{}, models.ErrConflict
		case 500:
			return models.ResponseBookISBN{}, models.ErrServerIssue
		default:
			return models.ResponseBookISBN{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var book models.ResponseBookISBN
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Printf("--ERROR-- with GetBookDetails(): %v\n", err)
		return models.ResponseBookISBN{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetBookDetails() OK")
	return book, nil
}

func (c *ExternalAPIClient) SearchForMovieByTitle(movieTitle string) (models.ResponseMovieSearch, error) {
	queryParameters := fmt.Sprintf("query=%s&include_adult=true", url.QueryEscape(movieTitle))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.SearchMovie, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForMovieByTitle(): %v\n", err)
		return models.ResponseMovieSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForMovieByTitle(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseMovieSearch{}, models.ErrBadRequest
		case 401:
			return models.ResponseMovieSearch{}, models.ErrUnauthorized
		case 404:
			return models.ResponseMovieSearch{}, models.ErrNotFound
		case 409:
			return models.ResponseMovieSearch{}, models.ErrConflict
		case 500:
			return models.ResponseMovieSearch{}, models.ErrServerIssue
		default:
			return models.ResponseMovieSearch{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var movie models.ResponseMovieSearch
	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Printf("--ERROR-- with SearchForMovieByTitle(): %v\n", err)
		return models.ResponseMovieSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForMovieByTitle() OK")
	return movie, nil
}

type parametersGetMoviesTvDetails struct {
	MovieID  string `json:"movie_id"`
	TvID     string `json:"tv_id"`
	Language string `json:"language"`
}

func (c *ExternalAPIClient) GetMovieDetails(movieID string) (models.ResponseMovieDetails, error) {

	params := parametersGetMoviesTvDetails{
		MovieID: movieID,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.GetDetails, params)
	if err != nil {
		log.Printf("--ERROR-- with GetMovieDetails(): %v\n", err)
		return models.ResponseMovieDetails{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetMovieDetails(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseMovieDetails{}, models.ErrBadRequest
		case 401:
			return models.ResponseMovieDetails{}, models.ErrUnauthorized
		case 404:
			return models.ResponseMovieDetails{}, models.ErrNotFound
		case 409:
			return models.ResponseMovieDetails{}, models.ErrConflict
		case 500:
			return models.ResponseMovieDetails{}, models.ErrServerIssue
		default:
			return models.ResponseMovieDetails{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var movieDetails models.ResponseMovieDetails
	err = json.NewDecoder(r.Body).Decode(&movieDetails)
	if err != nil {
		log.Printf("--ERROR-- with GetMovieDetails(): %v\n", err)
		return models.ResponseMovieDetails{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetMovieDetails() OK")
	return movieDetails, nil
}

func (c *ExternalAPIClient) GetMovieCredits(movieID string) (models.ResponseMovieCredits, error) {
	queryParameters := fmt.Sprintf("movie_id=%s", url.QueryEscape(movieID))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.GetMovieCredits, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetMovieCredits(): %v\n", err)
		return models.ResponseMovieCredits{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetMovieCredits(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseMovieCredits{}, models.ErrBadRequest
		case 401:
			return models.ResponseMovieCredits{}, models.ErrUnauthorized
		case 404:
			return models.ResponseMovieCredits{}, models.ErrNotFound
		case 409:
			return models.ResponseMovieCredits{}, models.ErrConflict
		case 500:
			return models.ResponseMovieCredits{}, models.ErrServerIssue
		default:
			return models.ResponseMovieCredits{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var movieCredits models.ResponseMovieCredits
	err = json.NewDecoder(r.Body).Decode(&movieCredits)
	if err != nil {
		log.Printf("--ERROR-- with GetMovieCredits(): %v\n", err)
		return models.ResponseMovieCredits{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetMovieCredits() OK")
	return movieCredits, nil
}

func (c *ExternalAPIClient) SearchForSeriesByTitle(seriesTitle string) (models.ResponseTvSearch, error) {
	queryParameters := fmt.Sprintf("query=%s&include_adult=true", url.QueryEscape(seriesTitle))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.SearchTV, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForSeriesByTitle(): %v\n", err)
		return models.ResponseTvSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForSeriesByTitle(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseTvSearch{}, models.ErrBadRequest
		case 401:
			return models.ResponseTvSearch{}, models.ErrUnauthorized
		case 404:
			return models.ResponseTvSearch{}, models.ErrNotFound
		case 409:
			return models.ResponseTvSearch{}, models.ErrConflict
		case 500:
			return models.ResponseTvSearch{}, models.ErrServerIssue
		default:
			return models.ResponseTvSearch{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var series models.ResponseTvSearch
	err = json.NewDecoder(r.Body).Decode(&series)
	if err != nil {
		log.Printf("--ERROR-- with SearchForSeriesByTitle(): %v\n", err)
		return models.ResponseTvSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForSeriesByTitle() OK")
	return series, nil
}

func (c *ExternalAPIClient) GetSeriesDetails(seriesID string) (models.ResponseTvDetails, error) {

	params := parametersGetMoviesTvDetails{
		TvID: seriesID,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.GetDetails, params)
	if err != nil {
		log.Printf("--ERROR-- with GetSeriesDetails(): %v\n", err)
		return models.ResponseTvDetails{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetSeriesDetails(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseTvDetails{}, models.ErrBadRequest
		case 401:
			return models.ResponseTvDetails{}, models.ErrUnauthorized
		case 404:
			return models.ResponseTvDetails{}, models.ErrNotFound
		case 409:
			return models.ResponseTvDetails{}, models.ErrConflict
		case 500:
			return models.ResponseTvDetails{}, models.ErrServerIssue
		default:
			return models.ResponseTvDetails{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var seriesDetails models.ResponseTvDetails
	err = json.NewDecoder(r.Body).Decode(&seriesDetails)
	if err != nil {
		log.Printf("--ERROR-- with GetSeriesDetails(): %v\n", err)
		return models.ResponseTvDetails{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetSeriesDetails() OK")
	return seriesDetails, nil
}

func (c *ExternalAPIClient) SearchForVideogameOnPlatformByTitle(videogameTitle, platform string) (models.ResponseVideogameSearch, error) {
	// Get right platform ID (based on RAWG)
	var platformID string
	platform = strings.ToLower(platform)
	switch platform {
	case "xbox one":
		platformID = "1"
	case "ios":
		platformID = "3"
	case "pc":
		platformID = "4"
	case "macos":
		platformID = "5"
	case "linux":
		platformID = "6"
	case "nintendo switch":
		platformID = "7"
	case "android":
		platformID = "21"
	case "playstation 4":
		platformID = "18"
	case "xbox series":
		platformID = "186"
	case "playstation 5":
		platformID = "185"
	default:
		platformID = ""
	}

	// Create query parameters
	queryParameters := ""
	if platformID != "" {
		queryParameters = fmt.Sprintf("search=%s&platforms=%s", url.QueryEscape(videogameTitle), platformID)
	} else {
		queryParameters = fmt.Sprintf("search=%s", url.QueryEscape(videogameTitle))
	}

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Videogames.Search, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForVideogameOnPlatformByTitle(): %v\n", err)
		return models.ResponseVideogameSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForVideogameOnPlatformByTitle(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseVideogameSearch{}, models.ErrBadRequest
		case 401:
			return models.ResponseVideogameSearch{}, models.ErrUnauthorized
		case 404:
			return models.ResponseVideogameSearch{}, models.ErrNotFound
		case 409:
			return models.ResponseVideogameSearch{}, models.ErrConflict
		case 500:
			return models.ResponseVideogameSearch{}, models.ErrServerIssue
		default:
			return models.ResponseVideogameSearch{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var videogame models.ResponseVideogameSearch
	err = json.NewDecoder(r.Body).Decode(&videogame)
	if err != nil {
		log.Printf("--ERROR-- with SearchForVideogameOnPlatformByTitle(): %v\n", err)
		return models.ResponseVideogameSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForVideogameOnPlatformByTitle() OK")
	return videogame, nil
}

func (c *ExternalAPIClient) GetVideogameDetails(videogameID string) (models.ResponseVideogameDetails, error) {
	queryParameters := fmt.Sprintf("id=%s", url.QueryEscape(videogameID))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Videogames.GetDetails, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetVideogameDetails(): %v\n", err)
		return models.ResponseVideogameDetails{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetVideogameDetails(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseVideogameDetails{}, models.ErrBadRequest
		case 401:
			return models.ResponseVideogameDetails{}, models.ErrUnauthorized
		case 404:
			return models.ResponseVideogameDetails{}, models.ErrNotFound
		case 409:
			return models.ResponseVideogameDetails{}, models.ErrConflict
		case 500:
			return models.ResponseVideogameDetails{}, models.ErrServerIssue
		default:
			return models.ResponseVideogameDetails{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var vgDetails models.ResponseVideogameDetails
	err = json.NewDecoder(r.Body).Decode(&vgDetails)
	if err != nil {
		log.Printf("--ERROR-- with GetVideogameDetails(): %v\n", err)
		return models.ResponseVideogameDetails{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetVideogameDetails() OK")
	return vgDetails, nil
}

func (c *ExternalAPIClient) SearchForBoardgameByTitle(boardgameTitle string) (models.ResponseBoardgameSearch, error) {
	queryParameters := fmt.Sprintf("query=%s", url.QueryEscape(boardgameTitle))

	// Make request to get boardgames list
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Boardgames.Search, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBoardgameByTitle(): %v\n", err)
		return models.ResponseBoardgameSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForBoardgameByTitle(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseBoardgameSearch{}, models.ErrBadRequest
		case 401:
			return models.ResponseBoardgameSearch{}, models.ErrUnauthorized
		case 404:
			return models.ResponseBoardgameSearch{}, models.ErrNotFound
		case 409:
			return models.ResponseBoardgameSearch{}, models.ErrConflict
		case 500:
			return models.ResponseBoardgameSearch{}, models.ErrServerIssue
		default:
			return models.ResponseBoardgameSearch{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var boardgames models.ResponseBoardgameSearch
	err = json.NewDecoder(r.Body).Decode(&boardgames)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBoardgameByTitle(): %v\n", err)
		return models.ResponseBoardgameSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForBoardgameByTitle() OK")
	return boardgames, nil
}

func (c *ExternalAPIClient) GetBoardgameDetails(boardgameID string) (models.ResponseBoardgameDetails, error) {
	queryParameters := fmt.Sprintf("id=%s", url.QueryEscape(boardgameID))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Boardgames.GetDetails, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with GetBoardgameDetails(): %v\n", err)
		return models.ResponseBoardgameDetails{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetBoardgameDetails(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ResponseBoardgameDetails{}, models.ErrBadRequest
		case 401:
			return models.ResponseBoardgameDetails{}, models.ErrUnauthorized
		case 404:
			return models.ResponseBoardgameDetails{}, models.ErrNotFound
		case 409:
			return models.ResponseBoardgameDetails{}, models.ErrConflict
		case 500:
			return models.ResponseBoardgameDetails{}, models.ErrServerIssue
		default:
			return models.ResponseBoardgameDetails{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var bgDetails models.ResponseBoardgameDetails
	err = json.NewDecoder(r.Body).Decode(&bgDetails)
	if err != nil {
		log.Printf("--ERROR-- with GetBoardgameDetails(): %v\n", err)
		return models.ResponseBoardgameDetails{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetBoardgameDetails() OK")
	return bgDetails, nil
}
