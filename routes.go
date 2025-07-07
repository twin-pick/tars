package tars

import (
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerRoutes() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/users/:usernames", s.findCommonFilm)
}

func (s *Server) findCommonFilm(c *gin.Context) {
	log.Info(time.Now())

	usernamesQuery := c.Param("usernames")
	usernames := strings.Split(usernamesQuery, ",")

	result, err := fetchScrapper(usernames)
	if err != nil {
		log.Errorf("Error fetching scrapper: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	if result == (Film{}) {
		log.Warn("No common films found")
		c.JSON(http.StatusNotFound, Header{"message": "No common films found"})
		return
	}

	log.Info(time.Now())

	log.Infof("Common film found: %s (%s)", result.Title, result.Date)
	c.JSON(http.StatusOK, result)
}
