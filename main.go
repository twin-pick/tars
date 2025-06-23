package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type Film struct {
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type tmdbSearchResponse struct {
	Results []struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		ID          int    `json:"id"`
	} `json:"results"`
}

type WatchList struct {
	Films []string `json:"films"`
}

func fetchTmdbFilm(title string) (Film, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?query=%s", title)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Film{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJlNDBmNDcwNGRiYjZhZGJlOWVhMTFmNmMyOTQxNmQ2ZiIsIm5iZiI6MTc1MDY3MDE2MC43NTEwMDAyLCJzdWIiOiI2ODU5MWI1MDJmMWQwNzg0MTQ0YmQ1NWUiLCJzY29wZXMiOlsiYXBpX3JlYWQiXSwidmVyc2lvbiI6MX0.d9dHVIaSdu5Q27P5pYCub4D369dZoWNN5U7kEhaVY6w")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Film{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Film{}, err
	}

	var searchResp tmdbSearchResponse
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
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch watchlists"})
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

			url := fmt.Sprintf("http://localhost:8000/api/user/watchlist/%s", u)

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

			var films []string
			if err := json.Unmarshal(body, &films); err != nil {
				log.Errorf("Error during JSON parsing for user %s: %v", u, err)
				resultChan <- WatchList{}
				return
			}

			resultChan <- WatchList{Films: films}
		}(username)
	}

	wg.Wait()
	close(resultChan)

	watchlists := []WatchList{}
	for wl := range resultChan {
		watchlists = append(watchlists, wl)
	}

	return compareAndFindCommonFilms(watchlists)
}

func compareAndFindCommonFilms(watchlists []WatchList) (Film, error) {
	if len(watchlists) == 0 {
		return Film{}, fmt.Errorf("No watchlists provided")}
	}

	var commonFilms []Film

	for _, film := range watchlists[0].Films {
		existsInAll := true

		for _, wl := range watchlists[1:] {
			if !watchlistContainsFilm(film, wl) {
				existsInAll = false
				break
			}
		}

		if existsInAll {
			film := Film{Title: film}
			commonFilms = append(commonFilms, film)
		}
	}

	return chooseRandomFilm(commonFilms)
}

func watchlistContainsFilm(film string, watchlist WatchList) bool {
	return slices.Contains(watchlist.Films, film)
}

func chooseRandomFilm(films []Film) (Film, error) {
	if len(films) == 0 {
		return Film{}, fmt.Errorf("No common films found")}
	}
	randNum := rand.Intn(len(films))
	return fetchTmdbFilm(films[randNum].Title)
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
