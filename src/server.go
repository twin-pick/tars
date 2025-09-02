package src

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewServer(cfg *Config) *Server {
	return &Server{
		Router: gin.Default(),
		Config: cfg,
		Rooms:  make(map[string]*Room),
	}
}

func (s *Server) Run() {
	s.registerRoutes()
	s.Router.Run(fmt.Sprintf(":%s", s.Config.ExposedPort))
}

func (s *Server) registerRoutes() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")

	v1.GET("/match/:usernames", s.match)
	v1.GET("/match/:usernames/:genres", s.match)
	v2.GET("/party/:usernames", s.party)
	v2.GET("/party/:usernames/:genres", s.party)
	v2.GET("/party/room/:roomId", s.room)
}
