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

type Duration string

type QueryParams struct {
	Usernames []string
	Genres    []string
	Duration  string
}

type Film struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Date     string `json:"date"`
	Poster   string `json:"poster"`
	Duration string `json:"duration"`
}

type WatchList struct {
	Films []*Film `json:"films"`
}

type OMDbResponse struct {
	Director string `json:"Director"`
	Poster   string `json:"Poster"`
	Duration string `json:"Runtime"`
}

type Client struct {
	Connection *WebsocketConn
	Votes      map[string]*Vote
}

type Room struct {
	Id        string
	Clients   map[string]*WebsocketConn
	Watchlist *WatchList
	Votes     map[string]*Vote
	Mutex     Mutex
}

type Vote struct {
	FilmId      string `json:"filmId"`
	WantToWatch bool   `json:"wantToWatch"`
	SocketId    string `json:"socketId"`
}

type EventIdentification struct {
	Event    string `json:"event"`
	SocketId string `json:"socketId"`
}

type EventData struct {
	Event string  `json:"event"`
	Data  []*Film `json:"films"`
}

type Result struct {
	Film *Film `json:"film"`
	Votes string `json:"votes"`
}

type VotesResults map[int]*Film

type EventResults struct {
	Event   string           `json:"event"`
	Results []*Result `json:"results"`
}

type EventFilmSelected struct {
	Event        string `json:"event"`
	FilmSelected *Film  `json:"film"`
}
