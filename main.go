package main

import (
	"log"

	"SiYue_Blog/internal/di"
)

func main() {
	server, err := di.InitHTTPServer()
	if err != nil {
		log.Fatalf("initialize server failed: %v", err)
	}

	log.Printf("blog backend starting at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
