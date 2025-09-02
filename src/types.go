package src

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Context = gin.Context
type Header = gin.H
type Router = gin.Engine
type WebsocketConn = websocket.Conn
type WebsocketUpgrader = websocket.Upgrader
type Mutex = sync.Mutex
type Request = http.Request
type WaitGroup = sync.WaitGroup

type Server struct {
	Router *Router
	Config *Config
	Rooms  map[string]*Room
	Mutex  Mutex
}

type Config struct {
	OMDbApiKey   string
	ScrapperPort string
	ExposedPort  string
}

type QueryParams struct {
	Usernames []string
	Genres    []string
}

type Film struct {
	Id       string `json:"id"`
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

type WebsocketClient struct {
	Id     string
	Client *WebsocketConn
}

type Room struct {
	Id        string
	Clients   map[*WebsocketConn]bool
	Watchlist *WatchList
	Votes     []*Vote
	Mutex     Mutex
}

type Vote struct {
	FilmId      string `json:"filmId"`
	WantToWatch bool   `json:"wantToWatch"`
	SocketId    string `json:"socketId"`
}
