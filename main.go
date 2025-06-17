package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type LetterboxdWatchlistFilm struct {
	Title     string `json:"title"`
	Year      int    `json:"year"`
	Author    string `json:"author"`
	PosterUrl string `json:"poster_url"`
}

type LetteboxdWatchlist struct {
	Films []LetterboxdWatchlistFilm `json:"films"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/twin", func(context *gin.Context) {
		context.String(http.StatusOK, "Pick!")
	})

	router.GET("/users/:usernames", func(context *gin.Context) {
		usernamesQuery := context.Param("usernames")
		usernames := strings.Split(usernamesQuery, ",")

		results := fetchScrapper(usernames)

		context.JSON(http.StatusOK, results)
	})

	return router
}

func fetchScrapper(usernames []string) []LetterboxdWatchlistFilm {
	var wg sync.WaitGroup
	resultChan := make(chan LetteboxdWatchlist, len(usernames))

	for _, username := range usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			url := fmt.Sprintf("http://localhost:8000/api/watchlist/%s", u)
			resp, err := http.Get(url)
			if err != nil {
				log.Errorf("Error for user %s: %v", u, err)
				resultChan <- LetteboxdWatchlist{}
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("Error during body reading for user %s: %v", u, err)
				resultChan <- LetteboxdWatchlist{}
				return
			}

			var data []LetterboxdWatchlistFilm
			if err := json.Unmarshal(body, &data); err != nil {
				log.Errorf("Error during JSON parsing for user %s: %v", u, err)
				resultChan <- LetteboxdWatchlist{}
				return
			}

			for i := range data {
				data[i].Author = u
			}

			resultChan <- LetteboxdWatchlist{Films: data}
		}(username)
	}

	wg.Wait()
	close(resultChan)

	allWatchlists := []LetteboxdWatchlist{}
	for wl := range resultChan {
		allWatchlists = append(allWatchlists, wl)
	}

	return compareAndFindCommonFilms(allWatchlists)
}

func compareAndFindCommonFilms(watchlists []LetteboxdWatchlist) []LetterboxdWatchlistFilm {
	if len(watchlists) == 0 {
		return nil
	}

	commonFilms := make(map[string]LetterboxdWatchlistFilm)
	for _, film := range watchlists[0].Films {
		key := fmt.Sprintf("%s-%d", film.Title, film.Year)
		commonFilms[key] = film
	}

	for _, watchlist := range watchlists[1:] {
		currentFilms := make(map[string]LetterboxdWatchlistFilm)
		for _, film := range watchlist.Films {
			key := fmt.Sprintf("%s-%d", film.Title, film.Year)
			if _, exists := commonFilms[key]; exists {
				currentFilms[key] = film
			}
		}
		commonFilms = currentFilms
	}

	results := make([]LetterboxdWatchlistFilm, 0, len(commonFilms))
	for _, film := range commonFilms {
		results = append(results, film)
	}

	return results
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
