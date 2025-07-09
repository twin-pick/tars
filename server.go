package tars

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewServer(cfg *Config) *Server {
	server := &Server{
		Router: gin.Default(),
		Config: cfg,
	}
	server.registerRoutes()
	return server
}

func (s *Server) Run() {
	s.Router.Run(fmt.Sprintf(":%s", s.Config.ExposedPort))
}

func (s *Server) registerRoutes() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/users/:usernames", s.findCommonFilm)
	v1.GET("/users/:usernames/:genres", s.findCommonFilm)
}
