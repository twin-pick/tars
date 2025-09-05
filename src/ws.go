package src

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var upgrader = WebsocketUpgrader{
	CheckOrigin: func(r *Request) bool { return true },
}

func NewRoom(watchlist *WatchList) *Room {
	return &Room{
		Id:        uuid.New().String(),
		Clients:   make(map[string]*WebsocketConn),
		Watchlist: watchlist,
		Votes:     make(map[string]*Vote),
	}
}

func NewResult(film *Film, votes int) *Result {
	return &Result{
		Film:  film,
		Votes: strconv.Itoa(votes),
	}
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

	room.sendInitialData(conn)

	go func() {
		defer room.handleDisconnect(conn)

		for {
			var vote *Vote
			if err := conn.ReadJSON(&vote); err != nil {
				log.Errorf("error reading vote: %v", err)
				break
			}

			log.Infof("Received vote in room %s: %+v", room.Id, vote)

			room.handleVote(vote)
		}
	}()
}

func (r *Room) connectToRoom(c *Context) *WebsocketConn {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("failed to upgrade websocket: %v", err)
		return nil
	}

	socketId := uuid.NewString()

	r.Mutex.Lock()
	r.Clients[socketId] = conn
	r.Mutex.Unlock()

	event := EventIdentification{
		Event:    "identification",
		SocketId: socketId,
	}

	if err := conn.WriteJSON(event); err != nil {
		log.Errorf("failed to send identification event: %v", err)
		conn.Close()
		r.removeClient(conn)
	}

	log.Infof("Client %s joined room %s", socketId, r.Id)
	return conn
}

func (r *Room) sendInitialData(conn *WebsocketConn) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	event := EventData{
		Event: "data",
		Data:  r.Watchlist.Films,
	}

	if err := conn.WriteJSON(event); err != nil {
		log.Errorf("failed to send data event: %v", err)
		conn.Close()
		r.removeClient(conn)
	}

	return nil
}

func (r *Room) handleDisconnect(conn *WebsocketConn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for sid, c := range r.Clients {
		if c == conn {
			delete(r.Clients, sid)
			break
		}
	}

	conn.Close()
	log.Infof("Client left room %s", r.Id)
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

func (r *Room) handleVote(vote *Vote) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Votes[vote.SocketId+vote.FilmId] = vote

	log.Infof("Vote recorded by socket %s: filmId: %s, wantToWatch: %t", vote.SocketId, vote.FilmId, vote.WantToWatch)
	log.Infof("Current votes in room %s", r.Id)

	for _, v := range r.Votes {
		log.Infof(" - socket %s: filmId: %s, wantToWatch: %t", v.SocketId, v.FilmId, v.WantToWatch)
	}

	if len(r.Votes) == len(r.Clients)*len(r.Watchlist.Films) {
		results := r.getRoomResults()
		r.broadcastVoteResults(results)
		log.Infof("All votes received in room %s", r.Id)
	}

	if selected := r.checkFilmSelected(); selected != nil {
		log.Infof("Film selected: %s", selected.Title)
		r.broadcastFilmSelected(selected)
		r.Votes = make(map[string]*Vote)
	}
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

func (r *Room) broadcastFilmSelected(film *Film) {
	event := EventFilmSelected{
		Event:        "film_selected",
		FilmSelected: film,
	}

	log.Infof("Broadcasting film_selected to %d clients: %+v", len(r.Clients), film)

	for _, client := range r.Clients {
		if err := client.WriteJSON(event); err != nil {
			log.Errorf("failed to send film_selected event: %v", err)
			client.Close()
			r.removeClient(client)
		}
	}
}

func (r *Room) removeClient(conn *WebsocketConn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for sid, c := range r.Clients {
		if c == conn {
			delete(r.Clients, sid)
			break
		}
	}
}

func (r *Room) getRoomResults() []*Result {
	results := []*Result{}

	for _, film := range r.Watchlist.Films {
		var voteCount int

		for _, vote := range r.Votes {
			if film.Id == vote.FilmId && vote.WantToWatch {
				voteCount++
			}
		}

		results = append(results, NewResult(film, voteCount))
	}

	return results
}

func (r *Room) broadcastVoteResults(results []*Result) {
	event := EventResults{
		Event:   "results",
		Results: results,
	}

	for _, client := range r.Clients {
		if err := client.WriteJSON(event); err != nil {
			log.Errorf("failed to send result event: %v", err)
			client.Close()
			r.removeClient(client)
		}
	}
}
