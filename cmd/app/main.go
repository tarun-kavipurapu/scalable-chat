package main

import (
	"log"
	"tarun-kavipurapu/test-go-chat/config"
	"tarun-kavipurapu/test-go-chat/internal"
)

func main() {
	config.LoadConfig(".")
	server := internal.NewHTTPServer()

	err := server.Start(":8080")
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

}
