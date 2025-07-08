package tars

func (s *Server) registerRoutes() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/users/:usernames", s.findCommonFilm)
	v1.GET("/users/:usernames/:genres", s.findCommonFilm)
}
