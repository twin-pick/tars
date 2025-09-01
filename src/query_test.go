package src

import "testing"

func TestGetUsernamesParams(t *testing.T) {
	usernames, err := getUsernamesParams("user1,user2,user3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(usernames) != 3 {
		t.Errorf("Expected 3 usernames, got %d", len(usernames))
	}
	if usernames[0] != "user1" || usernames[1] != "user2" || usernames[2] != "user3" {
		t.Errorf("Usernames do not match expected values")
	}
}

func TestGetGenresParams(t *testing.T) {
	genres := getGenresParams("war,action,adventure")
	if len(genres) != 3 {
		t.Errorf("Expected 3 genres, got %d", len(genres))
	}
	if genres[0] != "war" || genres[1] != "action" || genres[2] != "adventure" {
		t.Errorf("Genres do not match expected values")
	}
}
