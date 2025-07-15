package src

import (
	"net/http"

	"github.com/charmbracelet/log"
)

func (s *Server) findPartyFilm(c *Context) {
	querryParams := NewQueryParams(c)

	commonFilms, err := s.getCommonsFilms(querryParams)
	if err != nil {
		log.Errorf("Error getCommonsFilms: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	if len(commonFilms) == 0 {
		log.Warnf("No common films found")
		c.JSON(http.StatusOK, Header{"message": "No common films found"})
		return
	}

	commonFilmsWithDetails, err := s.getCommonsFilmsDetails(commonFilms)
	if err != nil {
		log.Errorf("Error getting common films details: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	log.Infof("Found %d common films", len(commonFilmsWithDetails))
	c.JSON(http.StatusOK, commonFilmsWithDetails)
}

func (s *Server) getCommonsFilmsDetails(films []*Film) ([]*Film, error) {
	for _, film := range films {
		film, err := s.getFilmDetails(film)
		if err != nil {
			log.Errorf("Error getting film details for %s: %v", film.Title, err)
			return nil, err
		}
	}
	return films, nil
}
