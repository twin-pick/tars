package src

import "strings"

func NewQueryParams(c *Context) *QueryParams {
	usernamesQuery := c.Param("usernames")
	usernames := strings.Split(usernamesQuery, ",")

	genresQuery := c.Param("genres")
	var genres []string
	if genresQuery != "" {
		genres = strings.Split(genresQuery, ",")
	}

	return &QueryParams{
		Usernames: usernames,
		Genres:    genres,
	}
}
