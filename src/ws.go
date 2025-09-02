package src

import (
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var upgrader = WebsocketUpgrader{
	CheckOrigin: func(r *Request) bool { return true },
}

func (s *Server) createRoom(room *Room) {
	s.Mutex.Lock()
	s.Rooms[room.Id] = room
	s.Mutex.Unlock()
	log.Infof("Created room %s with %d films", room.Id, len(room.Watchlist.Films))
}

func (s *Server) room(c *Context) {
	room, err := s.getRoom(c)
	if err != nil {
		return
	}

	conn := room.connectToRoom(c)
	if conn == nil {
		c.JSON(http.StatusInternalServerError, Header{"error": "failed to connect to room"})
		return
	}

	if err := room.sendInitialData(conn); err != nil {
		log.Errorf("failed to send initial data: %v", err)
		conn.Close()
		return
	}

	go func() {
		defer room.handleDisconnect(conn)

		for {
			var vote Vote
			if err := conn.ReadJSON(&vote); err != nil {
				log.Errorf("error reading vote: %v", err)
				break
			}

			log.Infof("Received vote in room %s: %+v", room.Id, vote)

			selected := room.handleVote(vote)

			if selected != nil {
				log.Infof("Film selected in room %s: %+v", room.Id, selected)
				room.broadcastFilmSelected(selected)
			}
		}
	}()
}

func (s *Server) getRoom(c *Context) (*Room, error) {
	roomID := c.Param("roomId")

	s.Mutex.Lock()
	room, ok := s.Rooms[roomID]
	s.Mutex.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, Header{"error": "room not found"})
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (r *Room) connectToRoom(c *Context) *WebsocketConn {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("failed to upgrade websocket: %v", err)
		return nil
	}

	socketId := uuid.NewString()

	r.Mutex.Lock()
	r.Clients[conn] = true
	r.Mutex.Unlock()

	conn.WriteJSON(map[string]string{
		"socketId": socketId,
	})

	log.Infof("Client %s joined room %s", socketId, r.Id)
	return conn
}

func (r *Room) sendInitialData(conn *WebsocketConn) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	return conn.WriteJSON(r.Watchlist.Films)
}

func (r *Room) handleDisconnect(conn *WebsocketConn) {
	r.Mutex.Lock()
	delete(r.Clients, conn)
	r.Mutex.Unlock()
	conn.Close()
	log.Infof("Client left room %s", r.Id)
}

func (r *Room) handleVote(vote Vote) *Film {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Votes = append(r.Votes, &vote)
	return r.checkFilmSelected()
}

func (r *Room) broadcastFilmSelected(film *Film) {
	r.Mutex.Lock()
	clientsCopy := make([]*WebsocketConn, 0, len(r.Clients))
	for client := range r.Clients {
		clientsCopy = append(clientsCopy, client)
	}
	r.Mutex.Unlock()

	event := map[string]any{
		"event": "film_selected",
		"film":  film,
	}

	for _, client := range clientsCopy {
		if err := client.WriteJSON(event); err != nil {
			log.Errorf("failed to send film_selected event: %v", err)
			client.Close()
			r.removeClient(client)
		}
	}
}

func (r *Room) removeClient(conn *WebsocketConn) {
	r.Mutex.Lock()
	delete(r.Clients, conn)
	r.Mutex.Unlock()
}

func (r *Room) checkFilmSelected() *Film {
	voteCount := make(map[string]int)
	for _, v := range r.Votes {
		if v.WantToWatch {
			voteCount[v.FilmId]++
		}
	}

	clientCount := len(r.Clients)
	for filmId, count := range voteCount {
		if count == clientCount {
			for _, f := range r.Watchlist.Films {
				if f.Id == filmId {
					return f
				}
			}
		}
	}
	return nil
}
