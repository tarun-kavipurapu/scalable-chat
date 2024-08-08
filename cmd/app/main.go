package main

import (
	"flag"
	"log"
	"os"
	"tarun-kavipurapu/test-go-chat/config"
	"tarun-kavipurapu/test-go-chat/internal"
)

func main() {
	// Set up command-line flags
	port := flag.String("port", "8080", "Port to run the HTTP server on")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.LoadConfig(".")

	// Initialize the HTTP server
	server := internal.NewHTTPServer()

	// Start the server on the specified port
	address := ":" + *port
	err := server.Start(address)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
