package main

import (
	. "github.com/twin-pick/tars/src"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	server := NewServer(config)
	server.RegisterRoutes()
	server.Run()
}
