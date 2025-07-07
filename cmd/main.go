package main

import (
	tars "github.com/twin-pick/tars"
)

func main() {
	config := tars.NewConfig()
	server := tars.NewServer(config)
	server.Run()
}
