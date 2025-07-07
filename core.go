package tars

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

func fetchScrapper(usernames []string) (Film, error) {
	var wg sync.WaitGroup
	resultChan := make(chan WatchList, len(usernames))

	log.Info(time.Now())

	for _, username := range usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			watchlist, err := fetchWatchlist(u)
			if err != nil {
				log.Errorf("Failed to fetch watchlist for user %s: %v", u, err)
				resultChan <- WatchList{}
				return
			}
			resultChan <- watchlist
		}(username)
	}

	log.Info(time.Now())

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

func fetchWatchlist(username string) (WatchList, error) {
	log.Infof("Fetching watchlist for user: %s", username)

	url := fmt.Sprintf("http://localhost:8000/api/v2/%s/watchlist", username)

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
