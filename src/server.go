package src

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *Config) *Server {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:8082"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return &Server{
		Router: router,
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

	v1.GET("/match", s.match)
	v2.GET("/party", s.party)
	v2.GET("/party/room/:roomId", s.room)
}
