package main

import "sync"

const NumWorkers = 5
const NumDownloaderWorkers = 2
const RefreshTokenDurationInMonth = 2
const youtubeBaseUrl = "https://www.youtube.com/watch"
const duplicateUserError = "pq: duplicate key value violates unique constraint"
const UnauthorizedErrorMessage = "this user is not authorized to make this request"

var BaseAlbumPath = "broke_music_app"

var AuthTypes = []string{"SELF", "GOOGLE"}

var TaskChannel = make(chan YtDlpTask)
var Tasks []YtDlpTask
var DownloadTasks []YtDownloadTask
var DownloadTaskChannel = make(chan YtDownloadTask)
var mutex sync.Mutex
