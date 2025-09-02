package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func NewFilm(title string, date string) *Film {
	return &Film{
		Title: title,
		Date:  date,
	}
}

func NewWatchlist(films []*Film) *WatchList {
	if films == nil {
		films = []*Film{}
	}
	for i := range films {
		films[i].Id = uuid.New().String()
	}
	return &WatchList{Films: films}
}

func (s *Server) fetchCommonsFilms(qp *QueryParams) ([]*Film, error) {
	var wg WaitGroup
	resultChan := make(chan *WatchList, len(qp.Usernames))

	for _, username := range qp.Usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			watchlist, err := s.fetchUserWatchlist(username, qp)
			if err != nil {
				log.Errorf("Failed to fetch watchlist for user %s: %v", u, err)
				resultChan <- &WatchList{}
				return
			}
			resultChan <- watchlist
		}(username)
	}

	wg.Wait()
	close(resultChan)

	var watchlists []*WatchList
	watchlists = make([]*WatchList, 0, len(resultChan))
	for wl := range resultChan {
		if len(wl.Films) != 0 {
			watchlists = append(watchlists, wl)
		}
	}

	commonFilms, err := s.compareWatchlists(watchlists)
	if err != nil {
		log.Errorf("Error comparing watchlists: %v", err)
		return []*Film{}, fmt.Errorf("error comparing watchlists: %w", err)
	}

	return commonFilms, nil
}

func (s *Server) fetchUserWatchlist(username string, qp *QueryParams) (*WatchList, error) {
	url := fmt.Sprintf("http://localhost:%s/api/v4/%s/watchlist", s.Config.ScrapperPort, username)

	if qp.Genres != nil {
		url += fmt.Sprintf("/%s", strings.Join(qp.Genres, ","))
	}

	resp, err := http.Get(url)
	if err != nil {
		return &WatchList{}, fmt.Errorf("error for user %s: %w", username, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &WatchList{}, fmt.Errorf("error reading body for user %s: %w", username, err)
	}

	var films []*Film
	if err := json.Unmarshal(body, &films); err != nil {
		return &WatchList{}, fmt.Errorf("error parsing JSON for user %s: %w", username, err)
	}

	log.Infof("Fetched %d films for user: %s", len(films), username)
	return NewWatchlist(films), nil
}

func (s *Server) compareWatchlists(watchlists []*WatchList) ([]*Film, error) {
	if len(watchlists) == 0 {
		return []*Film{}, fmt.Errorf("No watchlists provided")
	}

	filmCount := make(map[string]*Film)
	occurrences := make(map[string]int)

	for _, wl := range watchlists {
		seen := make(map[string]bool)
		for i := range wl.Films {
			film := wl.Films[i]
			if !seen[film.Title] {
				occurrences[film.Title]++
				if _, exists := filmCount[film.Title]; !exists {
					filmCount[film.Title] = film
				}
				seen[film.Title] = true
			}
		}
	}

	var commonFilms []*Film
	for title, count := range occurrences {
		if count == len(watchlists) {
			commonFilms = append(commonFilms, filmCount[title])
		}
	}

	return commonFilms, nil
}

func (s *Server) getFilmDetails(film *Film) (*Film, error) {
	escapedTitle := url.QueryEscape(film.Title)
	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&y=%s&apikey=%s", escapedTitle, film.Date, s.Config.OMDbApiKey)

	log.Infof("Fetching details for film: %s from URL: %s", film.Title, url)

	resp, err := http.Get(url)
	if err != nil {
		return &Film{}, fmt.Errorf("error fetching film details: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Film{}, fmt.Errorf("error reading response body: %w", err)
	}

	var omdbResponse OMDbResponse
	err = json.Unmarshal(body, &omdbResponse)
	if err != nil {
		return &Film{}, fmt.Errorf("error parsing OMDb response: %w", err)
	}

	film.Director = omdbResponse.Director
	film.Poster = omdbResponse.Poster

	return film, nil
}

func (s *Server) getCommonsFilms(c *Context) ([]*Film, error) {
	queryParams, err := NewQueryParams(c)
	if err != nil {
		log.Errorf("Error parsing query params: %v", err)
		c.JSON(http.StatusBadRequest, Header{"error": err.Error()})
		return nil, err
	}

	commonFilms, err := s.fetchCommonsFilms(queryParams)
	if err != nil {
		log.Errorf("Error getCommonsFilms: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return nil, err
	}

	if len(commonFilms) == 0 {
		log.Warnf("No common films found")
		c.JSON(http.StatusOK, Header{"message": "No common films found"})
		return nil, err
	}

	return commonFilms, nil
}
