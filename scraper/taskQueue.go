package scraper

import (
	"sync"
	"time"
)

type Task struct {
	URL string
}

type TaskQueue struct {
	mutex       sync.Mutex
	paused      bool
	processFunc func(Task)
	tasks       chan Task
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	queue := &TaskQueue{
		tasks:  make(chan Task, bufferSize),
		paused: false,
	}
	go queue.run()
	return queue
}

func (q *TaskQueue) AddTask(url string) {
	q.tasks <- Task{URL: url}
}

func (q *TaskQueue) run() {
	for task := range q.tasks {
		q.mutex.Lock()
		for q.paused {
			q.mutex.Unlock()
			time.Sleep(time.Second)
			q.mutex.Lock()
		}
		q.mutex.Unlock()
		q.processFunc(task)
	}
}

func (q *TaskQueue) SetProcessFunc(f func(Task)) {
	q.mutex.Lock()
	q.processFunc = f
	q.mutex.Unlock()
}

func (q *TaskQueue) Pause(t time.Duration) {
	q.mutex.Lock()
	q.paused = true
	q.mutex.Unlock()

	time.Sleep(t)

	q.mutex.Lock()
	q.paused = false
	q.mutex.Unlock()
}
