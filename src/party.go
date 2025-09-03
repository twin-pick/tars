package src

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func NewRoom(watchlist *WatchList) *Room {
	return &Room{
		Id:        uuid.New().String(),
		Clients:   make(map[*WebsocketConn]bool),
		Watchlist: watchlist,
	}
}

func (s *Server) party(c *Context) {
	qp, err := NewQueryParams(c)
	if err != nil {
		log.Errorf("Error parsing query params: %v", err)
		c.JSON(http.StatusBadRequest, Header{"error": err.Error()})
		return
	}

	commonFilms, err := s.getCommonsFilms(qp)
	if err != nil {
		log.Errorf("Error retrieving common films: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	if len(commonFilms) == 0 {
		log.Warnf("No common films found")
		c.JSON(http.StatusOK, Header{"message": "No common films found"})
		return
	}

	commonFilmsWithDetails, err := s.getCommonFilmsDetails(commonFilms)
	if err != nil {
		log.Errorf("Error getting common films details: %v", err)
		c.JSON(http.StatusInternalServerError, Header{"error": err.Error()})
		return
	}

	watchlist := NewWatchlist(commonFilmsWithDetails)
	room := NewRoom(watchlist)

	s.createRoom(room)

	log.Infof("Found %d common films", len(watchlist.Films))
	c.JSON(http.StatusOK, Header{
		"roomId": room.Id,
	})
}
