# brokeMusicApp
    BrokeMusicApp is an app which lets its users stream music online for free.
    This repo is for the server part of the app.

![github test badge](https://github.com/RazafimanantsoaJohnson/brokemusicapp/actions/workflows/ci.yml/badge.svg)

## Why this name?
    I just thought of a name for a music streaming app for broke people so there it was XD.

## Why use it when spotify already exists?
    Well, free Spotify don't let its users fully control their listening; whether you don't get to play the song you're interested in right away or the next song in your listening queue is completely random.
    It is particularly annoying when you're trying to listen to an album (where the artist wrote songs in a specific order to trigger specific feelings or tell a story).

## How it works:
    There was 2 big parts to resolving this problem:
        - How to get albums' data with the tracks they contain and their order
        - Where to find the actula albums' audio tracks
    1. For the albums' data ; I registered for spotify API to authenticate and getaccess to those data.
    2. For the audio tracks, I registered for youtube API (to search for the youtube url of each track in the album); and use *yt-dlp* to get the audio streaming url for the clients to be able to consume those streaming url right away. On the background the audios are downloaded in the server to build up the *'local music bank'* and use less and less youtube streaming url (avoiding issues like expiration or permission blocks at terms).

## Available endpoints:
    - **POST /api/signup**: create a new user
        Body (json):
            `
            {
                email: "email",
                password: "password"
            }`
    - **POST /api/signin**: create a new user
        Body (json):
            `
            {
                email: "email",
                password: "password"
            }`
    ### **Endpoint requiring a JWT (provided as a *Bearer token*)**
    - **GET /api/albums?query=XXX**: searching for albums using the 'XXX' keyword as a search parameter
    - **GET /api/users/favorites**: getting 10 recently visited albums by the current user ('visited' here means the user checked the album's tracks)
    - **GET /api/albums/{albumId}/tracks**: getting the list of tracks from an album with the specified id
    - **GET /api/albums/{albumId}/tracks/{trackId}**: getting specific track data based on ID (can be used to get the youtube streaming url of a track or refresh 'outdated data' on client side).
        **GET /api/albums/{albumId}/tracks/{trackId}?retry=true**: forcing the server to get a new streaming url for the track (will check if it the retry is valid by checking if the file is already downloaded or if the youtube url not responding with an 'OK' status)
    - **GET /api/albums/{albumId}/tracks/{trackId}/stream**: streams an 'already downloaded' audio file directly from the server. (will stream the audio binary)

## How to install:
    - PREREQUISITE softwares: **Go, PostgreSQL, Goose** 
    - The running computer needs to have **yt-dlp** installed and available as a running executable.
    - Environment variable (maybe provided by a **'.env'** file):
        DB_URL (url of the database connection)
        PORT (the port the server will listen to)
        SPOTIFY_CLIENTID (ID of our spotify registered app to enable OAuth)
        SPOTIFY_CLIENTSECRET (Secret of our spotify registered app to enable OAuth)
        YOUTUBE_APIKEY (Apikey we get from registering for youtube API)
    - Execute:
    `
        $go get .
        $ goose postgres "postgres://username:password@localhost:port/brokemusicapp" up
    `

## How to run:
    `
        $go run .
    `
        *or*
    `
        $go build
        $./brokemusicapp
    `

## Notes:
    To prevent the server from crashing, background work, and ensure concurrency on users requests for track's streaming urls, the app has a set number of *'dedicated worker go routines'* to interact with yt-dlp, set the found value in DB and if necessary, return the found value to the user request.

    ![userRequest_taskQueue_worker illustration](https://github.com/RazafimanantsoaJohnson/brokemusicapp/resources/request_taskQueue_worker_illustration.png)
