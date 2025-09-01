package src

import "testing"

func TestNewFilm(T *testing.T) {
	film := NewFilm("Inception", "2010")
	if film.Title != "Inception" {
		T.Errorf("Expected title 'Inception', got '%s'", film.Title)
	}
}

func TestNewWatchlist(T *testing.T) {
	watchlist := NewWatchlist(nil)
	if len(watchlist.Films) != 0 {
		T.Errorf("Expected empty watchlist, got %d films", len(watchlist.Films))
	}

	films := []*Film{NewFilm("Inception", "2010")}
	watchlist = NewWatchlist(films)
	if len(watchlist.Films) != 1 {
		T.Errorf("Expected watchlist with 1 film, got %d films", len(watchlist.Films))
	}
}
