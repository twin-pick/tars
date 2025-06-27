package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var tmdbToken string
var scrapperPort string
var port string

type FilmEntry struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

type Film struct {
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type TMDBSearchResponse struct {
	Results []struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		ID          int    `json:"id"`
	} `json:"results"`
}

type WatchList struct {
	Films []Film `json:"films"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/twin", func(context *gin.Context) {
		context.String(http.StatusOK, "Pick!")
	})

	router.GET("/users/:usernames", func(context *gin.Context) {
		usernamesQuery := context.Param("usernames")
		usernames := strings.Split(usernamesQuery, ",")
		result, err := fetchScrapper(usernames)
		if err != nil {
			log.Errorf("Error fetching scrapper: %v", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if result == (Film{}) {
			log.Info("No common films found")
			context.JSON(http.StatusNotFound, gin.H{"message": "No common films found"})
			return
		}

		context.JSON(http.StatusOK, result)
	})

	return router
}

func fetchScrapper(usernames []string) (Film, error) {
	var wg sync.WaitGroup
	resultChan := make(chan WatchList, len(usernames))

	for _, username := range usernames {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			url := fmt.Sprintf("http://localhost:8000/api/v1/%s/watchlist", u)

			resp, err := http.Get(url)
			if err != nil {
				log.Errorf("Error for user %s: %v", u, err)
				resultChan <- WatchList{}
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("Error during body reading for user %s: %v", u, err)
				resultChan <- WatchList{}
				return
			}
			// log.Infof("Response body for user %s: %s", u, string(body))

			var entries []FilmEntry
			if err := json.Unmarshal(body, &entries); err != nil {
				log.Errorf("Error during JSON parsing for user %s: %v", u, err)
				resultChan <- WatchList{}
				return
			}
			films := make([]Film, len(entries))
			for i, e := range entries {
				year := 0
				if len(e.Date) >= 4 {
					year, _ = strconv.Atoi(e.Date[:4])
				}
				films[i] = Film{
					Title:    e.Title,
					Year:     year,
					Director: "",
				}
			}

			resultChan <- WatchList{Films: films}
		}(username)
	}

	wg.Wait()
	close(resultChan)

	watchlists := []WatchList{}
	for wl := range resultChan {
		if len(wl.Films) != 0 {
			watchlists = append(watchlists, wl)
		}
	}

	return compareAndFindCommonFilms(watchlists)
}

func compareAndFindCommonFilms(watchlists []WatchList) (Film, error) {
	if len(watchlists) == 0 {
		return Film{}, fmt.Errorf("No watchlists provided")
	}

	var commonFilms []Film

	for _, film := range watchlists[0].Films {
		existsInAll := true

		for _, wl := range watchlists[1:] {
			if !watchlistContainsFilm(film.Title, wl) { // on compare par titre uniquement
				existsInAll = false
				break
			}
		}

		if existsInAll {
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
	// return fetchTmdbFilm(films[randNum].Title)
	return films[randNum], nil
}

func fetchTmdbFilm(title string) (Film, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?query=%s", title)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Film{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tmdbToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Film{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Film{}, err
	}

	var searchResp TMDBSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return Film{}, err
	}

	if len(searchResp.Results) == 0 {
		return Film{}, fmt.Errorf("Any result found for '%s'", title)
	}

	first := searchResp.Results[0]

	year := 0
	if len(first.ReleaseDate) >= 4 {
		year, _ = strconv.Atoi(first.ReleaseDate[:4])
	}

	return Film{
		Title:    first.Title,
		Year:     year,
		Director: "",
	}, nil
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

func main() {
	loadEnv()

	port = fmt.Sprintf(":%s", port)
	router := setupRouter()
	router.Run(port)
}
