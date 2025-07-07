package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var tmdbToken string
var scrapperPort string
var port string

type Film struct {
	Title string `json:"title"`
	Date  int    `json:"date"`
}

type WatchList struct {
	Films []Film `json:"films"`
}

type Router = gin.Engine
type Context = gin.Context
type Header = gin.H

func fetchScrapper(usernames []string) (Film, error) {
	var wg sync.WaitGroup
	resultChan := make(chan WatchList, len(usernames))

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

	var commonFilms []Film

	for _, film := range watchlists[0].Films {
		existsInAll := true

		for _, wl := range watchlists[1:] {
			if !watchlistContainsFilm(film.Title, wl) {
				existsInAll = false
				break
			}
		}

		if existsInAll {
			log.Infof("Found common film: %s (%d)", film.Title, film.Date)
			commonFilms = append(commonFilms, film)
		}
	}

	return chooseRandomFilm(commonFilms)
}

func watchlistContainsFilm(title string, watchlist WatchList) bool {
	for _, film := range watchlist.Films {
		if film.Title == title {
			return true
		}
	}
	return false
}

func chooseRandomFilm(films []Film) (Film, error) {
	if len(films) == 0 {
		return Film{}, fmt.Errorf("No common films found")
	}
	randNum := rand.Intn(len(films))
	return films[randNum], nil
}

func loadEnv() {
	godotenv.Load()

	tmdbToken = os.Getenv("TMDB_TOKEN")
	if tmdbToken == "" {
		log.Fatal("TMDB_TOKEN env var is not set")
	}

	scrapperPort = os.Getenv("SCRAPPER_PORT")
	if scrapperPort == "" {
		log.Fatal("SCRAPPER_PORT env var is not set")
	}

	port = os.Getenv("EXPOSED_PORT")
	if port == "" {
		log.Fatal("EXPOSED_PORT env var is not set")
	}
}

func createRouter() *Router {
	router := gin.Default()

	router.GET("/users/:usernames", func(context *Context) {
		usernamesQuery := context.Param("usernames")
		usernames := strings.Split(usernamesQuery, ",")
		result, err := fetchScrapper(usernames)
		if err != nil {
			log.Errorf("Error fetching scrapper: %v", err)
			context.JSON(http.StatusInternalServerError, Header{"error": err})
			return
		}

		if result == (Film{}) {
			log.Info("No common films found")
			context.JSON(http.StatusNotFound, Header{"message": "No common films found"})
			return
		}

		log.Infof("Common film found: %s (%d)", result.Title, result.Date)
		context.JSON(http.StatusOK, result)
	})

	return router
}

func main() {
	loadEnv()

	router := createRouter()
	router.Run(fmt.Sprintf(":%s", port))
}
