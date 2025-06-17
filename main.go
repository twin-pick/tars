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

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/twin", func(context *gin.Context) {
		context.String(http.StatusOK, "Pick!")
	})

	router.GET("/users/:usernames", func(context *gin.Context) {
		usernamesQuery := context.Param("usernames")
		usernames := strings.Split(usernamesQuery, ";")

		results := fetchScrapper(usernames)

		context.JSON(http.StatusOK, results)
	})

	return router
}

func fetchScrapper(usernames []string) []LetterboxdWatchlistFilm {
	var wg sync.WaitGroup
	resultChan := make(chan []LetterboxdWatchlistFilm, len(usernames))

	for _, username := range usernames {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			url := fmt.Sprintf("http://localhost:8000/api/watchlist/%s", u)
			resp, err := http.Get(url)
			if err != nil {
				log.Error("Error for user %s: %v\n", u, err)
				resultChan <- []LetterboxdWatchlistFilm{}
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error during body reading for user %s: %v\n", u, err)
				resultChan <- []LetterboxdWatchlistFilm{}
				return
			}

			var data []LetterboxdWatchlistFilm
			if err := json.Unmarshal(body, &data); err != nil {
				log.Error("Error during JSON parsing for user %s: %v\n", u, err)
				resultChan <- []LetterboxdWatchlistFilm{}
				return
			}

			resultChan <- data
		}(username)
	}

	wg.Wait()
	close(resultChan)

	finalResults := []LetterboxdWatchlistFilm{}
	for r := range resultChan {
		finalResults = append(finalResults, r...)
	}

	return finalResults
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
