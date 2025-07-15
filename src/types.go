package src

import "github.com/gin-gonic/gin"

type Context = gin.Context
type Header = gin.H
type Router = gin.Engine

type Server struct {
	Router *Router
	Config *Config
}

type Config struct {
	OMDBApiKey   string
	ScrapperPort string
	ExposedPort  string
}

type QueryParams struct {
	Usernames []string
	Genres    []string
}

type Film struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Date     string `json:"date"`
	Poster   string `json:"poster"`
}

type WatchList struct {
	Films []*Film `json:"films"`
}

type OMDbResponse struct {
	Director string `json:"Director"`
	Poster   string `json:"Poster"`
}

type Vote struct {
	Film         *Film
	WantsToWatch bool
}
