package tars

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"

	"github.com/charmbracelet/log"
)

func (s *Server) findCommonFilm(c *Context) {
	queryParams := NewQueryParams(c)

	result, err := s.fetchScrapper(queryParams)
	if err != nil {
		log.Errorf("Error fetchScrapper: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	if isEmptyFilm(result) {
		log.Warn("No common films found")
		c.JSON(http.StatusNotFound, Header{"message": "No common films found"})
		return
	}

	log.Infof("Common film found: %s (%s)", result.Title, result.Date)
	c.JSON(http.StatusOK, result)
}

func isEmptyFilm(film Film) bool {
	return film == Film{}
}

func (s *Server) fetchScrapper(qp *QueryParams) (Film, error) {
	var wg sync.WaitGroup
	resultChan := make(chan WatchList, len(qp.Usernames))

	for _, username := range qp.Usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			watchlist, err := s.fetchUserWatchlist(qp)
			if err != nil {
				log.Errorf("Failed to fetch watchlist for user %s: %v", u, err)
				resultChan <- WatchList{}
				return
			}
			resultChan <- watchlist
		}(username)
	}

	wg.Wait()
	close(resultChan)

	var watchlists []WatchList
	watchlists = make([]WatchList, 0, len(resultChan))
	for wl := range resultChan {
		if len(wl.Films) != 0 {
			watchlists = append(watchlists, wl)
		}
	}

	return s.compareAndFindCommonFilm(watchlists)
}

func (s *Server) fetchUserWatchlist(qp *QueryParams) (WatchList, error) {
	url := fmt.Sprintf("http://localhost:%s/api/v2/%s/watchlist", s.Config.ScrapperPort, qp.Usernames)

	for _, genre := range qp.Genres {
		log.Infof("Adding genre filter: %s", genre)
	}

	resp, err := http.Get(url)
	if err != nil {
		return WatchList{}, fmt.Errorf("error for user %s: %w", qp.Usernames, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WatchList{}, fmt.Errorf("error reading body for user %s: %w", qp.Usernames, err)
	}

	var films []Film
	if err := json.Unmarshal(body, &films); err != nil {
		return WatchList{}, fmt.Errorf("error parsing JSON for user %s: %w", qp.Usernames, err)
	}

	log.Infof("Fetched %d films for user: %s", len(films), qp.Usernames)

	return WatchList{Films: films}, nil
}

func (s *Server) compareAndFindCommonFilm(watchlists []WatchList) (Film, error) {
	if len(watchlists) == 0 {
		return Film{}, fmt.Errorf("No watchlists provided")
	}

	filmCount := make(map[string]*Film)
	occurrences := make(map[string]int)

	for _, wl := range watchlists {
		seen := make(map[string]bool)
		for i := range wl.Films {
			film := &wl.Films[i]
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

	return s.selectRandomFilm(commonFilms)
}

func (s *Server) selectRandomFilm(films []*Film) (Film, error) {
	if len(films) == 0 {
		return Film{}, fmt.Errorf("No common films found")
	}
	randNum := rand.Intn(len(films))
	return s.getFilmDetails(films[randNum])
}

func (s *Server) getFilmDetails(film *Film) (Film, error) {
	escapedTitle := url.QueryEscape(film.Title)
	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&y=%s&apikey=%s", escapedTitle, film.Date, s.Config.OMDBApiKey)

	log.Infof("Fetching details for film: %s from URL: %s", film.Title, url)

	resp, err := http.Get(url)
	if err != nil {
		return Film{}, fmt.Errorf("error fetching film details: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Film{}, fmt.Errorf("error reading response body: %w", err)
	}

	var omdb OMDbResponse
	err = json.Unmarshal(body, &omdb)
	if err != nil {
		return Film{}, fmt.Errorf("error parsing OMDb response: %w", err)
	}

	film.Director = omdb.Director
	film.Poster = omdb.Poster

	return *film, nil
}
