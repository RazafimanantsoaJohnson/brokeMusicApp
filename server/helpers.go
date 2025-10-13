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

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/auth"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/youtube"
	"github.com/google/uuid"
)

type YtDlpTask struct {
	YoutubeId  string
	TrackId    uuid.UUID // the id in DB
	Priority   int       // should just be 0,1  (or 1,2)
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

func downloadFile(url string) (bool, error) { // we probably don't want to see the errors
	tmpFile, err := os.Create("tracks_tmp/test.webm")
	if err != nil {
		return false, err
	}
	defer tmpFile.Close()
	testUrl := `https://rr3---sn-h50gpup0nuxaxjvh-hg0d.googlevideo.com/videoplayback?expire=1760099866&ei=uqnoaJ7GCZWXsvQPj8eRqAU&ip=102.117.235.41&id=o-AERuKwKJRRUKvqo62JGC8a-PghYFIb8U2ghPijhqWc1g&itag=251&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&met=1760078266%2C&mh=S2&mm=31%2C29&mn=sn-h50gpup0nuxaxjvh-hg0d%2Csn-hc57enee&ms=au%2Crdu&mv=m&mvi=3&pl=21&rms=au%2Cau&initcwndbps=2227500&bui=ATw7iSUMp6lzlWKoC3yVNz_z0qwswxOr25xfDe2wgQeY6dmpGJxPNUNn-Hil6AHvPv13aER_ymYN3ADn&vprv=1&svpuc=1&mime=audio%2Fwebm&ns=JSn06xFfeWmjEMBug9rUND0Q&rqh=1&gir=yes&clen=5045095&dur=308.721&lmt=1739333989535137&mt=1760077755&fvip=5&keepalive=yes&lmw=1&fexp=51557447%2C51565116%2C51565681%2C51580970&c=TVHTML5&sefc=1&txp=4532534&n=ClSoVRuS8uDFJQ&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cbui%2Cvprv%2Csvpuc%2Cmime%2Cns%2Crqh%2Cgir%2Cclen%2Cdur%2Clmt&lsparams=met%2Cmh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Crms%2Cinitcwndbps&lsig=APaTxxMwRgIhAM8vM4kljHkEuP1z_Inb3gbGsNOMuNKujr9WfGTNN77lAiEAnAyXFD9Z9Rh4kk2DOJi8r6VjkNt-30pac15E9SQWHKI%3D&sig=AJfQdSswRQIhANc90GdSRlb417uwfBn149LY9xyMMEz6TGHq-HZrxYO0AiBVeevVvSKutfPUITTCUzIVxX5jaM5PRTvzeBepq1xgkw%3D%3D`
	response, err := http.Get(testUrl)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, err
	}
	io.Copy(tmpFile, response.Body)
	return true, nil
}

func StartWorkerPool(cfg *ApiConfig) {
	for i := 0; i < NumWorkers; i++ {
		go worker(i, cfg)
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

func worker(id int, cfg *ApiConfig) {
	for task := range TaskChannel {
		// we are here supposing we are only passing one video per worker
		fmt.Printf("Worker %v is treating video (%v)\n", id, task.YoutubeId)
		videoUrl := fmt.Sprintf("%v?v=%v", youtubeBaseUrl, task.YoutubeId)
		result := YtDlpTaskResult{}
		urlParam := []string{videoUrl}
		extractedJson, err := youtube.CallYtDlpCmd(urlParam)
		if err != nil {
			result.err = err
		}
		result.result = extractedJson[0]
		audioStreamingUrl := (youtube.GetAudioStreamingUrl(extractedJson[0])).Url
		if task.ResultChan == nil {
			cfg.db.InsertTrackYoutubeUrl(context.Background(), database.InsertTrackYoutubeUrlParams{
				ID:         task.TrackId,
				Youtubeurl: sql.NullString{String: audioStreamingUrl, Valid: true},
			})
			continue
		}
		task.ResultChan <- result
	}
}

func (cfg *ApiConfig) middlewareCheckAuth(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		receivedToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(UnauthorizedErrorMessage))
			return
		}
		_, err = auth.ValidateJWT(receivedToken, cfg.jwtSecret)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(UnauthorizedErrorMessage))
		}
		next(w, r)
	}
}
