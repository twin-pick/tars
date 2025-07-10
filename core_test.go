package tars

import (
	"slices"
	"testing"
)

func film(title string) Film {
	return Film{Title: title, Date: "2000"}
}

func (s *Server) TestCompareAndFindCommonFilms(t *testing.T) {
	tests := []struct {
		name       string
		watchlists []WatchList
		wantOneOf  []string
		expectErr  bool
	}{
		{
			name: "one common film",
			watchlists: []WatchList{
				{Films: []Film{film("A"), film("B")}},
				{Films: []Film{film("B"), film("C")}},
			},
			wantOneOf: []string{"B"},
		},
		{
			name: "multiple common films",
			watchlists: []WatchList{
				{Films: []Film{film("A"), film("B"), film("C")}},
				{Films: []Film{film("B"), film("C"), film("D")}},
				{Films: []Film{film("C"), film("B"), film("E")}},
			},
			wantOneOf: []string{"B", "C"},
		},
		{
			name: "no common films",
			watchlists: []WatchList{
				{Films: []Film{film("A")}},
				{Films: []Film{film("B")}},
			},
			expectErr: true,
		},
		{
			name:       "empty input",
			watchlists: []WatchList{},
			expectErr:  true,
		},
		{
			name: "watchlist with duplicate films",
			watchlists: []WatchList{
				{Films: []Film{film("X"), film("X"), film("Y")}},
				{Films: []Film{film("Y"), film("X")}},
			},
			wantOneOf: []string{"X", "Y"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := s.compareAndFindCommonFilm(tc.watchlists)

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			found := false
			found = slices.Contains(tc.wantOneOf, result.Title)

			if !found {
				t.Errorf("expected one of %v, got %s", tc.wantOneOf, result.Title)
			}
		})
	}
}
