package tars

func main() {
	config := NewConfig()
	server := NewServer(config)
	server.Run()
}
