package main

import "sync"

const NumWorkers = 5
const youtubeBaseUrl = "https://www.youtube.com/watch"
const duplicateUserError = "pq: duplicate key value violates unique constraint"
const UnauthorizedErrorMessage = "this user is not authorized to make this request"

var AuthTypes = []string{"SELF", "GOOGLE"}

var TaskChannel = make(chan YtDlpTask)
var Tasks []YtDlpTask
var mutex sync.Mutex
