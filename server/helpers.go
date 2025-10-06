package main

import (
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/youtube"
	"github.com/google/uuid"
)

type YtDlpTask struct {
	TrackId    uuid.UUID // the id in DB
	Priority   int       // should just be 0,1  (or 1,2)
	ResultChan chan youtube.YtDlpExtractedJson
}

// the normal case would be that we just pop at the beginning of the array (might change on priorities) => we'll check the value to be poped
func pushTask(tasks []YtDlpTask, newTask YtDlpTask) {
	// FIFO for all but give place to priorities
	if len(tasks) == 0 || newTask.Priority == 0 {
		tasks = append(tasks, newTask)
		return
	}
	tasksCp := tasks[:]
	// tasks will be locked while treating batches anyway so no worries
	indexOfInsertion := 0
	for i, v := range tasks {
		if i == len(tasksCp)-1 {
			tasks = append(tasks, newTask) // a slice of priority 1
			return
		}
		if v.Priority == 0 && i == 0 {
			tasks = append([]YtDlpTask{newTask}, tasksCp...) // put it first
			return
		}
		if v.Priority == 0 {
			indexOfInsertion = i
			highPriorities := tasksCp[:indexOfInsertion]
			lowPriorities := tasksCp[indexOfInsertion:]
			tasks = append(highPriorities, newTask)
			tasks = append(tasks, lowPriorities...)
		}
	}
}

func popTask(tasks []YtDlpTask, taskToRemove YtDlpTask) {
	indexOfTaskToRemove := 0
	for i, v := range tasks {
		if v.TrackId == taskToRemove.TrackId {
			indexOfTaskToRemove = i
		}
	}
	if indexOfTaskToRemove == 0 {

	}
	beforeTask := tasks[:indexOfTaskToRemove]
	afterTask := tasks[indexOfTaskToRemove+1:]
	tasks = append(beforeTask, afterTask...)
}

func returnJson() {

}
