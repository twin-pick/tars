package src

import (
	"errors"
	"strings"

	"github.com/charmbracelet/log"
)

func NewQueryParams(c *Context) (*QueryParams, error) {
	usernames, err := getUsernamesParams(c.Param("usernames"))
	if err != nil {
		return nil, err
	}
	genres := getGenresParams(c.Param("genres"))

	return &QueryParams{
		Usernames: usernames,
		Genres:    genres,
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
