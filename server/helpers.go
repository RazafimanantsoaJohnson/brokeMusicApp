package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/youtube"
	"github.com/google/uuid"
)

type YtDlpTask struct {
	YoutubeId  string
	AlbumId    string
	TrackId    uuid.UUID // the id in DB
	Priority   int       // should just be 0,1  (or 1,2)  // -1 (no task), 0 (low task), 1 (high task)
	ResultChan chan YtDlpTaskResult
}

type YtDownloadTask struct {
	AlbumId                string
	TrackId                string
	YoutubeStreamingFormat youtube.YtDlpFormat
}

type YtDlpTaskResult struct {
	result youtube.YtDlpExtractedJson
	err    error
}

// the normal case would be that we just pop at the beginning of the array (might change on priorities) => we'll check the value to be poped
func pushTask(tasks *[]YtDlpTask, newTask YtDlpTask) {
	// FIFO for all but give place to priorities
	tasksCp := (*tasks)[:]
	if len(tasksCp) == 0 || newTask.Priority == 0 {
		*tasks = append(tasksCp, newTask)
		return
	}
	// tasks will be locked while treating batches anyway so no worries
	indexOfInsertion := 0
	for i, v := range tasksCp {
		if i == len(tasksCp)-1 {
			*tasks = append(tasksCp, newTask) // a slice of priority 1
			return
		}
		if v.Priority == 0 && i == 0 {
			*tasks = append([]YtDlpTask{newTask}, tasksCp...) // put it first
			return
		}
		if v.Priority > 0 {
			indexOfInsertion = i
			highPriorities := tasksCp[:indexOfInsertion]
			lowPriorities := tasksCp[indexOfInsertion:]
			*tasks = append(highPriorities, newTask)
			*tasks = append((*tasks), lowPriorities...)
		}
	}
}

func pop[T interface{}](tasks *[]T) T { // the '0' task will always be the biggest priority
	cpTasks := (*tasks)
	if len(cpTasks) == 1 {
		task := cpTasks[0]
		cpTasks := []T{}
		*tasks = cpTasks
		return task
	}
	task := cpTasks[0]
	*tasks = cpTasks[1:]
	return task
}

func returnJson[T interface{}](w http.ResponseWriter, value T) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(jsonValue)
}

func downloadFile(cfg *ApiConfig) { // we probably don't want to see the errors
	// if we are unauthorized; run yt-dlp again
	for task := range DownloadTaskChannel {
		fileName := fmt.Sprintf("tracks_tmp/%v.%v", task.AlbumId, task.TrackId)
		tmpFile, _ := os.Create(fileName) // will change to createTemp
		// if err != nil {
		// 	return false, err
		// }
		defer tmpFile.Close()
		response, _ := http.Get(task.YoutubeStreamingFormat.Url)
		// if err != nil {
		// 	return false, err
		// }
		defer response.Body.Close()

		// if response.StatusCode != http.StatusOK {
		// 	return false, err
		// }
		io.Copy(tmpFile, response.Body)
		time.Sleep(10 * time.Second) // place 10 seconds of pause between 2 downloads for the same worker

	}
}

func StartWorkerPool(cfg *ApiConfig) {
	for i := 0; i < NumWorkers; i++ {
		go worker(i, cfg)
	}

	go scheduler()

	// for i := 0; i < NumDownloaderWorkers; i++ {
	// 	go downloadFile(cfg)
	// }
	// go downloadsScheduler()
}

func scheduler() {
	// will run forever waiting for videos to process
	for {
		mutex.Lock()
		if len(Tasks) > 0 {
			task := pop(&Tasks)
			mutex.Unlock()
			TaskChannel <- task
		} else {
			mutex.Unlock()
			time.Sleep(600 * time.Millisecond)
		}
	}
}

func downloadsScheduler() {
	for {
		mutex.Lock()
		if len(DownloadTasks) > 0 {
			task := pop(&DownloadTasks)
			mutex.Unlock()
			DownloadTaskChannel <- task
		} else {
			time.Sleep(10 * time.Second)
			continue
		}
	}
}

func worker(id int, cfg *ApiConfig) {
	for task := range TaskChannel {
		fmt.Printf("Worker %v is treating video (%v)\n", id, task.YoutubeId)
		videoUrl := fmt.Sprintf("%v?v=%v", youtubeBaseUrl, task.YoutubeId)
		result := YtDlpTaskResult{}
		urlParam := []string{videoUrl}
		extractedJson, err := youtube.CallYtDlpCmd(urlParam)
		if err != nil {
			result.err = err
		}
		result.result = extractedJson[0]
		audioStreamingFormat := youtube.GetAudioStreamingUrl(extractedJson[0])
		if task.ResultChan == nil {
			cfg.db.InsertTrackYoutubeUrl(context.Background(), database.InsertTrackYoutubeUrlParams{
				ID:         task.TrackId,
				Youtubeurl: sql.NullString{String: audioStreamingFormat.Url, Valid: true},
			})

			// DownloadTasks = append(DownloadTasks, YtDownloadTask{
			// 	YoutubeStreamingFormat: audioStreamingFormat,
			// 	AlbumId:                task.AlbumId,
			// 	TrackId:                task.TrackId.String(),
			// })

			continue
		}
		task.ResultChan <- result
	}
}
