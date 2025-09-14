package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	port                string
	spotifyClientId     string
	spotifyClientSecret string
	spotifyAccessToken  string
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	spotifyClientId := os.Getenv("SPOTIFY_CLIENTID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENTSECRET")

	if port == "" {
		log.Fatalf("No port number was provided for the server")
	}
	if spotifyClientSecret == "" || spotifyClientId == "" {
		log.Fatalf("No spotify credentials was provided")
	}

	config := ApiConfig{
		port:                port,
		spotifyClientId:     spotifyClientId,
		spotifyClientSecret: spotifyClientSecret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello handsum, thank you for doing this"))
	})

	server := &http.Server{
		Addr:    ":" + config.port,
		Handler: mux,
	}

	fmt.Printf("The server is running and listening to port %s\n", config.port)
	server.ListenAndServe()
}
