package src

import (
	"net/http"
	"sync"

	"github.com/charmbracelet/log"
)

func (s *Server) findPartyFilm(c *Context) {
	queryParams := NewQueryParams(c)

	commonFilms, err := s.getCommonsFilms(queryParams)
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
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		firstError error
	)

	for i, film := range films {
		wg.Add(1)
		go func(i int, film *Film) {
			defer wg.Done()
			details, err := s.getFilmDetails(film)
			if err != nil {
				log.Errorf("Error getting film details for %s: %v", film.Title, err)
				mu.Lock()
				if firstError == nil {
					firstError = err
				}
				mu.Unlock()
				return
			}
			mu.Lock()
			films[i] = details
			mu.Unlock()
		}(i, film)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	return films, nil
}
