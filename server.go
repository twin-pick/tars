package tars

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg Config) *Server {
	server := &Server{
		Router: gin.Default(),
		Config: cfg,
	}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.Router.GET("/users/:usernames", s.handleUserWatchlist)
}

func (s *Server) handleUserWatchlist(c *gin.Context) {
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

func (s *Server) Run() {
	s.Router.Run(fmt.Sprintf(":%s", s.Config.ExposedPort))
}
