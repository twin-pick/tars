package src

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/charmbracelet/log"
)

func (s *Server) findCommonFilm(c *Context) {
	queryParams := NewQueryParams(c)

	commonFilms, err := s.getCommonsFilms(queryParams)
	if err != nil {
		log.Errorf("Error getCommonsFilms: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	if len(commonFilms) == 0 {
		log.Warnf("No common films found")
		c.JSON(http.StatusNotFound, Header{"message": "No common films found"})
		return
	}

	film, err := s.selectRandomFilm(commonFilms)
	if err != nil {
		log.Errorf("Error selecting random film: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	log.Infof("Common film found: %s (%s)", film.Title, film.Date)
	c.JSON(http.StatusOK, film)
}

func (s *Server) selectRandomFilm(films []*Film) (*Film, error) {
	if len(films) == 0 {
		return &Film{}, fmt.Errorf("No common films found")
	}
	randNum := rand.Intn(len(films))
	return s.getFilmDetails(films[randNum])
}
