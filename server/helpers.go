package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/youtube"
)

type YtDlpTask struct {
	YoutubeId string
	// TrackId    uuid.UUID // the id in DB
	Priority   int // should just be 0,1  (or 1,2)
	ResultChan chan YtDlpTaskResult
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

func pop(tasks *[]YtDlpTask) YtDlpTask { // the '0' task will always be the biggest priority
	cpTasks := (*tasks)
	if len(cpTasks) == 1 {
		return cpTasks[0]
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

func StartWorkerPool() {
	for i := 0; i < NumWorkers; i++ {
		go worker(i)
	}

	go scheduler()

	go func() {
		// simulation
		receivedTaskChannel := <-TaskChannel
		fmt.Println(receivedTaskChannel)
		fmt.Println(Tasks)
	}()
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

func worker(id int) {
	for task := range TaskChannel {
		// we are here supposing we are only passing one video per worker
		fmt.Printf("Worker %v is treating video %v ()\n", id, task.YoutubeId)
		videoUrl := fmt.Sprintf("%v?v=%v", youtubeBaseUrl, task.YoutubeId)
		result := YtDlpTaskResult{}
		urlParam := []string{videoUrl}
		extractedJson, err := youtube.CallYtDlpCmd(urlParam)
		if err != nil {
			result.err = err
		}
		result.result = extractedJson[0]
		fmt.Printf("Worker %v finished treating video %v \n", id, task.YoutubeId)
		fmt.Println("Extracted Video :", result.result.Title, "\t(good job worker)")
		task.ResultChan <- result
	}
}
