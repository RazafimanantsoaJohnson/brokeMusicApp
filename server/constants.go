package main

import "sync"

const NumWorkers = 5
const NumDownloaderWorkers = 2
const youtubeBaseUrl = "https://www.youtube.com/watch"

var TaskChannel = make(chan YtDlpTask)
var Tasks []YtDlpTask
var DownloadTasks []YtDownloadTask
var DownloadTaskChannel = make(chan YtDownloadTask)
var mutex sync.Mutex
