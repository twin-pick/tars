package tars

import (
	"fmt"

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

func (s *Server) Run() {
	s.Router.Run(fmt.Sprintf(":%s", s.Config.ExposedPort))
}
