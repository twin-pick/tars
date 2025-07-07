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
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	TMDBToken    string
	ScrapperPort string
	ExposedPort  string
}

type Server struct {
	Router *Router
	Config Config
}

type Router = gin.Engine
type Context = gin.Context
type Header = gin.H

type Film struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

type WatchList struct {
	Films []Film `json:"films"`
}

func loadEnv() Config {
	godotenv.Load()

	cfg := Config{
		TMDBToken:    os.Getenv("TMDB_TOKEN"),
		ScrapperPort: os.Getenv("SCRAPPER_PORT"),
		ExposedPort:  os.Getenv("EXPOSED_PORT"),
	}

	if cfg.TMDBToken == "" {
		log.Fatal("TMDB_TOKEN env var is not set")
	}
	if cfg.ScrapperPort == "" {
		log.Fatal("SCRAPPER_PORT env var is not set")
	}
	if cfg.ExposedPort == "" {
		log.Fatal("EXPOSED_PORT env var is not set")
	}

	return cfg
}

func NewServer(cfg Config) *Server {
	server := &Server{
		Router: gin.Default(),
		Config: cfg,
	}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.Router.GET("/users/:usernames", s.handleUserWatchlist)
}

func (s *Server) handleUserWatchlist(c *gin.Context) {
	log.Info(time.Now())

	usernamesQuery := c.Param("usernames")
	usernames := strings.Split(usernamesQuery, ",")

	result, err := fetchScrapper(usernames)
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

	log.Info(time.Now())

	log.Infof("Common film found: %s (%s)", result.Title, result.Date)
	c.JSON(http.StatusOK, result)
}

func (s *Server) Run() {
	s.Router.Run(fmt.Sprintf(":%s", s.Config.ExposedPort))
}

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

	log.Info("Fetched %d films for user: %s", len(films), username)

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

func main() {
	config := loadEnv()
	server := NewServer(config)
	server.Run()
}
