package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/logging"
)

func checkYoutubeUrlResponse(trackUrl string) bool {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ytUrlResponse, err := http.Get(trackUrl)
	if err != nil {
		logging.LogData(err.Error())
		return false
	}
	waitTimeOver := <-ctxWithTimeout.Done()
	if ytUrlResponse == nil {
		return false
	}
	if ytUrlResponse.StatusCode > 299 {
		return false
	}
	fmt.Println(waitTimeOver)
	return true
}
