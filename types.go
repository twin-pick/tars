package tars

import "github.com/gin-gonic/gin"

type Router = gin.Engine
type Context = gin.Context
type Header = gin.H

type Config struct {
	TMDBToken    string
	ScrapperPort string
	ExposedPort  string
}

type Server struct {
	Router *Router
	Config Config
}

type QueryParams struct {
	genres []string
}

type Film struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

type WatchList struct {
	Films []Film `json:"films"`
}
