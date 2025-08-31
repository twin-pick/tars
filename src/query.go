package src

import "strings"

func NewQueryParams(c *Context) *QueryParams {
	usernames := getUsernamesParams(c.Param("usernames"))
	genres := getGenresParams(c.Param("genres"))

	return &QueryParams{
		Usernames: usernames,
		Genres:    genres,
	}
}

func getUsernamesParams(queryParams string) []string {
	usernames := strings.Split(queryParams, ",")
	return usernames
}

func getGenresParams(queryParams string) []string {
	genres := strings.Split(queryParams, ",")
	return genres
}
