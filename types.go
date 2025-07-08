package tars

import "github.com/gin-gonic/gin"

type Router = gin.Engine
type Context = gin.Context
type Header = gin.H

type Config struct {
	OMDBApiKey   string
	ScrapperPort string
	ExposedPort  string
}

type Server struct {
	Router *Router
	Config Config
}

type Film struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Date     string `json:"date"`
	Poster   string `json:"poster"`
}

type WatchList struct {
	Films []Film `json:"films"`
}

type OMDbResponse struct {
	Director string `json:"Director"`
	Poster   string `json:"Poster"`
}
