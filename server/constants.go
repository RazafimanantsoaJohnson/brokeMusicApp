package main

import "sync"

const NumWorkers = 5
const youtubeBaseUrl = "https://www.youtube.com/watch"

var TaskChannel = make(chan YtDlpTask)
var Tasks []YtDlpTask
var mutex sync.Mutex
