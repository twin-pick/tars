package src

import (
	"errors"
	"strings"

	"github.com/charmbracelet/log"
)

func NewQueryParams(c *Context) (*QueryParams, error) {
	usernames, err := getUsernamesParams(c.Query("usernames"))
	if err != nil {
		return nil, err
	}

	genres := getGenresParams(c.Query("genres"))
	duration := getDurationParams(c.Query("duration"))

	return &QueryParams{
		Usernames: usernames,
		Genres:    genres,
		Duration:  duration,
	}, nil
}

func getUsernamesParams(queryParams string) ([]string, error) {
	if queryParams == "" {
		log.Errorf("No usernames provided")
		return nil, errors.New("no usernames provided")
	}
	usernames := strings.Split(queryParams, ",")
	return usernames, nil
}

func getGenresParams(queryParams string) []string {
	if queryParams == "" {
		return []string{}
	}
	genres := strings.Split(queryParams, ",")
	return genres
}

func getDurationParams(queryParams string) string {
	if queryParams == "" {
		return ""
	}
	return queryParams
}
