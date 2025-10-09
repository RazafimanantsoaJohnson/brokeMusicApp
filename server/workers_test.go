package main

import (
	"fmt"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	StartWorkerPool()

	cases := []struct {
		id         string
		resultChan chan YtDlpTaskResult
	}{
		{id: "diIFhc_Kzng", resultChan: make(chan YtDlpTaskResult)}, {id: "AE005nZeF-A", resultChan: make(chan YtDlpTaskResult)},
		{id: "uzS3WG6__G4", resultChan: make(chan YtDlpTaskResult)}, {id: "X_SEwgDl02E", resultChan: make(chan YtDlpTaskResult)},
		{id: "r4l9bFqgMaQ", resultChan: make(chan YtDlpTaskResult)}, {id: "HWDaIRe8_XI", resultChan: make(chan YtDlpTaskResult)},
		{id: "ncqkC9Ob2ZI", resultChan: make(chan YtDlpTaskResult)}, {id: "Dlz_XHeUUis", resultChan: make(chan YtDlpTaskResult)},
		{id: "P18g4rKns6Q", resultChan: make(chan YtDlpTaskResult)},
	}

	for _, c := range cases {
		mutex.Lock()
		pushTask(&Tasks, YtDlpTask{
			YoutubeId:  c.id,
			Priority:   0,
			ResultChan: c.resultChan,
		})
		// fmt.Println(Tasks)
		mutex.Unlock()
		go func() {
			workerResult := <-c.resultChan
			fmt.Println("Extracted Video :", workerResult.result.Title, "\t(good job worker)")
			fmt.Println(Tasks)
		}()
	}
	// Addition of a 'priority task'
	mutex.Lock()
	pushTask(&Tasks, YtDlpTask{
		YoutubeId:  "2npegbvmfso",
		Priority:   1,
		ResultChan: make(chan YtDlpTaskResult),
	})
	// pushTask(&Tasks, YtDlpTask{
	// 	YoutubeId:  "DloZ1xZHCmo",
	// 	Priority:   1,
	// 	ResultChan: make(chan YtDlpTaskResult),
	// })

	mutex.Unlock()
	mutex.Lock()
	pushTask(&Tasks, YtDlpTask{
		YoutubeId:  "DloZ1xZHCmo",
		Priority:   1,
		ResultChan: make(chan YtDlpTaskResult),
	})

	mutex.Unlock()

	time.Sleep(3 * time.Minute)
}
