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
		log.Warn("No common films found")
		c.JSON(http.StatusNotFound, Header{"message": "No common films found"})
		return
	}

	log.Infof("Found %d common films", len(commonFilms))
	c.JSON(http.StatusOK, commonFilms)
}
