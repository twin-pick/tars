package main

import (
	. "github.com/twin-pick/tars/src"
)

func main() {
	config := NewConfig()
	server := NewServer(config)
	server.Run()
}
