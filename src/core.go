package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

func NewWatchlist(films []*Film) *WatchList {
	if films == nil {
		films = []*Film{}
	}
	return &WatchList{Films: films}
}

func (s *Server) getCommonsFilms(qp *QueryParams) ([]*Film, error) {
	var wg sync.WaitGroup
	resultChan := make(chan *WatchList, len(qp.Usernames))

	for _, username := range qp.Usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			watchlist, err := s.fetchUserWatchlist(username, qp)
			if err != nil {
				logError("Failed to fetch watchlist for user %s: %v", u, err)
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

	commonFilms, err := s.compareAndFindCommonFilms(watchlists)
	if err != nil {
		logError("Error comparing watchlists: %v", err)
		return []*Film{}, fmt.Errorf("error comparing watchlists: %w", err)
	}

	return commonFilms, nil
}

func (s *Server) fetchUserWatchlist(username string, qp *QueryParams) (*WatchList, error) {
	url := fmt.Sprintf("http://localhost:%s/api/v2/%s/watchlist", s.Config.ScrapperPort, username)

	for _, genre := range qp.Genres {
		logInfo("Adding genre filter: %s", genre)
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

	logInfo("Fetched %d films for user: %s", len(films), username)

	return NewWatchlist(films), nil
}

func (s *Server) compareAndFindCommonFilms(watchlists []*WatchList) ([]*Film, error) {
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
	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&y=%s&apikey=%s", escapedTitle, film.Date, s.Config.OMDBApiKey)

	logInfo("Fetching details for film: %s from URL: %s", film.Title, url)

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
