package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/spotify"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ApiConfig struct {
	port                string
	spotifyClientId     string
	spotifyClientSecret string
	ytApiKey            string
	jwtSecret           string
	spotifyAccessToken  spotify.AuthResponse //string
	db                  *database.Queries
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	spotifyClientId := os.Getenv("SPOTIFY_CLIENTID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENTSECRET")
	ytApiKey := os.Getenv("YOUTUBE_APIKEY")
	dbUrl := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if port == "" {
		log.Fatalf("No port number was provided for the server")
	}
	if spotifyClientSecret == "" || spotifyClientId == "" {
		log.Fatalf("No spotify credentials was provided")
	}

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("%v", err)
	}
	db := database.New(connection)

	config := ApiConfig{
		port:                port,
		spotifyClientId:     spotifyClientId,
		spotifyClientSecret: spotifyClientSecret,
		ytApiKey:            ytApiKey,
		db:                  db,
		jwtSecret:           jwtSecret,
	}

	StartWorkerPool(&config)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello handsum, thank you for doing this"))
	})
	mux.HandleFunc("/api/signup", config.HandleSignup)
	mux.HandleFunc("/api/signin", config.HandleSignin)
	mux.HandleFunc("/api/albums", config.middlewareCheckAuth(config.HandleSearchAlbum))
	mux.HandleFunc("/api/albums/{albumId}/tracks", config.middlewareCheckAuth(config.HandleGetAlbumTracks))
	mux.HandleFunc("/api/albums/{albumId}/tracks/{trackId}", config.middlewareCheckAuth(config.HandleGetTrack))

	server := &http.Server{
		Addr:    ":" + config.port,
		Handler: mux,
	}

	fmt.Printf("The server is running and listening to port %s\n", config.port)
	server.ListenAndServe()
}
