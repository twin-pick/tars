package src

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/charmbracelet/log"
)

func (s *Server) match(c *Context) {
	qp, err := NewQueryParams(c)
	if err != nil {
		log.Errorf("Error parsing query params: %v", err)
		c.JSON(http.StatusBadRequest, Header{"error": err.Error()})
		return
	}

	commonFilms, err := s.getCommonsFilms(qp)
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

	var filmSelected *Film

	if qp.Duration != "" {
		commonFilmsWithDetails, err := s.getCommonFilmsDetails(commonFilms)
		if err != nil {
			log.Errorf("Error getting film details: %v", err)
			c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
			return
		}

		duration := parseDuration(qp.Duration)
		filtered := filterFilmsByDuration(commonFilmsWithDetails, duration)

		if len(filtered) == 0 {
			c.JSON(http.StatusNotFound, Header{"message": "No films found for this duration"})
			return
		}

		filmSelected, err = s.selectRandomFilm(filtered)
		if err != nil {
			log.Errorf("Error selecting random film: %v", err)
			c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		}
	} else {
		filmSelected, err = s.selectRandomFilm(commonFilms)
		if err != nil {
			log.Errorf("Error selecting random film: %v", err)
			c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
			return
		}
	}

	log.Infof("Common film found: %s (%s)", filmSelected.Title, filmSelected.Date)
	c.JSON(http.StatusOK, filmSelected)
}

func (s *Server) selectRandomFilm(films []*Film) (*Film, error) {
	if len(films) == 0 {
		return &Film{}, fmt.Errorf("No common films found")
	}
	randNum := rand.Intn(len(films))
	return s.getFilmDetails(films[randNum])
}
