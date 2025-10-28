package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	mux.HandleFunc("POST /api/signup", config.HandleSignup)
	mux.HandleFunc("POST /api/signin", config.HandleSignin)
	mux.HandleFunc("GET /api/albums", config.middlewareCheckAuth(HandleSearchAlbum))
	mux.HandleFunc("GET /api/albums/{albumId}/tracks", config.middlewareCheckAuth(HandleGetAlbumTracks))
	mux.HandleFunc("GET /api/albums/{albumId}/tracks/{trackId}", config.middlewareCheckAuth(HandleGetTrack))
	mux.HandleFunc("GET /api/albums/{albumId}/tracks/{trackId}/stream", config.middlewareCheckAuth(HandleServeTrackFile))
	mux.HandleFunc("GET /api/users/favorites", config.middlewareCheckAuth(HandleGetUserVisitedAlbums))

	server := &http.Server{
		Addr:              ":" + config.port,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Minute,
	}

	fmt.Printf("The server is running and listening to port %s\n", config.port)
	log.Fatal(server.ListenAndServe())
}
