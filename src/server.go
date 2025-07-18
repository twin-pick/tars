package src

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewServer(cfg *Config) *Server {
	return &Server{
		Router: gin.Default(),
		Config: cfg,
	}
}

func (s *Server) Run() {
	s.Router.Run(fmt.Sprintf(":%s", s.Config.ExposedPort))
}

func (s *Server) RegisterRoutes() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")

	v1.GET("/common/:usernames", s.findCommonFilm)
	v1.GET("/common/:usernames/:genres", s.findCommonFilm)
	v2.GET("/party/:usernames", s.findPartyFilm)
}
