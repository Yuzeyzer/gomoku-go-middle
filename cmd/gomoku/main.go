package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yuzeyzer/gomoku/internal/web"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	srv := web.NewServer(15)

	log.Printf("Gomoku web UI: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, srv.Handler()))
}
