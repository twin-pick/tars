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
			log.Errorf("Error getCommonFilmsDetails: %v", err)
			c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
			return
		}

		commonFilmsFilteredByDuration := filterFilmsByDuration(commonFilmsWithDetails, qp.Duration)
		filmSelected, err = s.selectRandomFilm(commonFilmsFilteredByDuration)
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

		filmSelected, err = s.getFilmDetails(filmSelected)
		if err != nil {
			log.Errorf("Error getFilmDetails: %v", err)
			c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
			return
		}
	}

	log.Infof("Common film found: %s (%s)", filmSelected.Title, filmSelected.Date)
	c.JSON(http.StatusOK, filmSelected)
}

func (s *Server) selectRandomFilm(films []*Film) (*Film, error) {
	if len(films) == 0 {
		return &Film{}, fmt.Errorf("no films available to select from")
	}
	randNum := rand.Intn(len(films))
	return films[randNum], nil
}
