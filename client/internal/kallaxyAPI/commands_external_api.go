package kallaxyapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *ExternalAPIClient) SearchForMovie(movieTitle string) (models.ResponseMovieSearch, error) {
	queryParameters := fmt.Sprintf("query=%s", url.QueryEscape(movieTitle))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.SearchMovie, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForMovie(): %v\n", err)
		return models.ResponseMovieSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForMovie(). Response status code: %v\n", r.StatusCode)
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
		log.Printf("--ERROR-- with SearchForMovie(): %v\n", err)
		return models.ResponseMovieSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForMovie() OK")
	return movie, nil
}

func (c *ExternalAPIClient) SearchForSeries(seriesTitle string) (models.ResponseTvSearch, error) {
	queryParameters := fmt.Sprintf("query=%s", url.QueryEscape(seriesTitle))

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.MoviesTV.SearchTV, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForSeries(): %v\n", err)
		return models.ResponseTvSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForSeries(). Response status code: %v\n", r.StatusCode)
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
		log.Printf("--ERROR-- with SearchForSeries(): %v\n", err)
		return models.ResponseTvSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForSeries() OK")
	return series, nil
}

func (c *ExternalAPIClient) SearchForVideogameOnSteam(videogameTitle string) (models.ResponseVideogameSearch, error) {
	queryParameters := fmt.Sprintf("search=%s&stores=%s", url.QueryEscape(videogameTitle), "1")

	// Make request
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Videogames.Search, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForVideogameOnSteam(): %v\n", err)
		return models.ResponseVideogameSearch{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForVideogameOnSteam(). Response status code: %v\n", r.StatusCode)
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
		log.Printf("--ERROR-- with SearchForVideogameOnSteam(): %v\n", err)
		return models.ResponseVideogameSearch{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForVideogameOnSteam() OK")
	return videogame, nil
}

func (c *ExternalAPIClient) SearchForBoardgame(boardgameTitle string) (models.ResponseBoardgameDetails, error) {
	queryParameters := fmt.Sprintf("query=%s&exact=1", url.QueryEscape(boardgameTitle))

	// Make request to get boardgames list
	r, err := c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Boardgames.Search, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBoardgame(): %v\n", err)
		return models.ResponseBoardgameDetails{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForBoardgame(). Response status code: %v\n", r.StatusCode)
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
	var boardgames models.ResponseBoardgameSearchAlternative
	err = json.NewDecoder(r.Body).Decode(&boardgames)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBoardgame() - Part boardgameS: %v\n", err)
		return models.ResponseBoardgameDetails{}, err
	}

	boardgameID := boardgames.Items.Item.ID
	// Make request for first board game of result list
	queryParameters = fmt.Sprintf("id=%s", url.QueryEscape(boardgameID))
	r, err = c.apiClient.makeHttpRequestWithQueryParameters(c.apiClient.Config.Endpoints.ExternalAPI.Boardgames.GetDetails, queryParameters)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBoardgame(): %v\n", err)
		return models.ResponseBoardgameDetails{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with SearchForBoardgame(). Response status code: %v\n", r.StatusCode)
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
	var boardgame models.ResponseBoardgameDetails
	err = json.NewDecoder(r.Body).Decode(&boardgame)
	if err != nil {
		log.Printf("--ERROR-- with SearchForBoardgame() - Part Boardgame: %v\n", err)
		return models.ResponseBoardgameDetails{}, err
	}

	// Return data
	log.Println("--DEBUG-- SearchForBoardgame() OK")
	return boardgame, nil
}
