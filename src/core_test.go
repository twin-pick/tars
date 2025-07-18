package src

import (
	"testing"
)

func TestNewFilm(T *testing.T) {
	film := NewFilm("Inception", "2010-07-16")

	if film.Title != "Inception" {
		T.Errorf("Expected title 'Inception', got '%s'", film.Title)
	}
	if film.Date != "2010-07-16" {
		T.Errorf("Expected date '2010-07-16', got '%s'", film.Date)
	}
}

func TestNewWatchlist(T *testing.T) {
	film1 := NewFilm("Inception", "2010-07-16")
	film2 := NewFilm("The Matrix", "1999-03-31")
	watchlist := NewWatchlist([]*Film{film1, film2})

	if len(watchlist.Films) != 2 {
		T.Errorf("Expected 2 films in watchlist, got %d", len(watchlist.Films))
	}
	if watchlist.Films[0].Title != "Inception" {
		T.Errorf("Expected first film title 'Inception', got '%s'", watchlist.Films[0].Title)
	}
	if watchlist.Films[1].Title != "The Matrix" {
		T.Errorf("Expected second film title 'The Matrix', got '%s'", watchlist.Films[1].Title)
	}
	if watchlist.Films[0].Date != "2010-07-16" {
		T.Errorf("Expected first film date '2010-07-16', got '%s'", watchlist.Films[0].Date)
	}
	if watchlist.Films[1].Date != "1999-03-31" {
		T.Errorf("Expected second film date '1999-03-31', got '%s'", watchlist.Films[1].Date)
	}
	if watchlist.Films[0] != film1 {
		T.Errorf("Expected first film to be film1, got '%v'", watchlist.Films[0])
	}
}
