package src

import "testing"

func TestNewServer(T *testing.T) {
	cfg, _ := NewConfig()
	server := NewServer(cfg)

	if server == nil {
		T.Error("Expected server to be created, but got nil")
	}
	if server.Router == nil {
		T.Error("Expected server router to be initialized, but got nil")
	}
}
