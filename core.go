package tars

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

func (s *Server) findCommonFilm(c *Context) {
	usernamesQuery := c.Param("usernames")
	usernames := strings.Split(usernamesQuery, ",")

	genresQuery := c.Param("genres")
	var genres []string
	if genresQuery != "" {
		genres = strings.Split(genresQuery, ",")
	}

	result, err := fetchScrapper(usernames, genres)
	if err != nil {
		log.Errorf("Error fetching scrapper: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	if result == (Film{}) {
		log.Warn("No common films found")
		c.JSON(http.StatusNotFound, Header{"message": "No common films found"})
		return
	}

	log.Infof("Common film found: %s (%s)", result.Title, result.Date)
	c.JSON(http.StatusOK, result)
}

func fetchScrapper(usernames []string, genres []string) (Film, error) {
	var wg sync.WaitGroup
	resultChan := make(chan WatchList, len(usernames))

	for _, username := range usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			watchlist, err := fetchWatchlist(u, genres)
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
	for wl := range resultChan {
		if len(wl.Films) != 0 {
			watchlists = append(watchlists, wl)
		}
	}

	return compareAndFindCommonFilms(watchlists)
}

func fetchWatchlist(username string, genres []string) (WatchList, error) {
	url := fmt.Sprintf("http://localhost:8000/api/v2/%s/watchlist", username)

	for _, genre := range genres {
		log.Infof("Adding genre filter: %s", genre)
	}

	resp, err := http.Get(url)
	if err != nil {
		return WatchList{}, fmt.Errorf("error for user %s: %w", username, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WatchList{}, fmt.Errorf("error reading body for user %s: %w", username, err)
	}

	var films []Film
	if err := json.Unmarshal(body, &films); err != nil {
		return WatchList{}, fmt.Errorf("error parsing JSON for user %s: %w", username, err)
	}

	log.Infof("Fetched %d films for user: %s", len(films), username)

	return WatchList{Films: films}, nil
}

func compareAndFindCommonFilms(watchlists []WatchList) (Film, error) {
	if len(watchlists) == 0 {
		return Film{}, fmt.Errorf("No watchlists provided")
	}

	filmCount := make(map[string]Film)
	occurrences := make(map[string]int)

	for _, wl := range watchlists {
		seen := make(map[string]bool)
		for _, film := range wl.Films {
			if !seen[film.Title] {
				occurrences[film.Title]++
				if _, exists := filmCount[film.Title]; !exists {
					filmCount[film.Title] = film
				}
				seen[film.Title] = true
			}
		}
	}

	var commonFilms []Film
	for title, count := range occurrences {
		if count == len(watchlists) {
			commonFilms = append(commonFilms, filmCount[title])
		}
	}

	return chooseRandomFilm(commonFilms)
}

func chooseRandomFilm(films []Film) (Film, error) {
	if len(films) == 0 {
		return Film{}, fmt.Errorf("No common films found")
	}
	randNum := rand.Intn(len(films))
	return films[randNum], nil
}
