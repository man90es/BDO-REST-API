package scraper

import (
	"bdo-rest-api/utils"
	"slices"
	"sync"
	"time"
)

type Task struct {
	ClientIP string
	Hash     string
	URL      string
}

type TaskQueue struct {
	clientIPs   map[string]int
	hashes      []string
	mutex       sync.Mutex
	paused      bool
	processFunc func(Task)
	tasks       chan Task
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	queue := &TaskQueue{
		clientIPs: make(map[string]int),
		paused:    false,
		tasks:     make(chan Task, bufferSize),
	}
	go queue.run()
	return queue
}

func (q *TaskQueue) AddTask(clientIP, hash, url string) {
	fullURL := utils.BuildRequest(url, map[string]string{
		"taskClient": clientIP,
		"taskHash":   hash,
	})

	q.mutex.Lock()
	if duplicate := slices.Contains(q.hashes, hash); duplicate {
		q.mutex.Unlock()
		return
	}
	q.clientIPs[clientIP]++
	q.hashes = append(q.hashes, hash)
	q.mutex.Unlock()

	q.tasks <- Task{
		ClientIP: clientIP,
		Hash:     hash,
		URL:      fullURL,
	}
}

func (q *TaskQueue) run() {
	for task := range q.tasks {
		q.mutex.Lock()
		for q.paused {
			q.mutex.Unlock()
			// FIXME: This is probably inefficient af
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

func (q *TaskQueue) CountQueuedTasksForClient(clientIP string) (count int) {
	q.mutex.Lock()
	count = max(0, q.clientIPs[clientIP])
	q.mutex.Unlock()

	return
}

func (q *TaskQueue) ConfirmTaskCompletion(clientIP string, hash string) {
	q.mutex.Lock()
	q.clientIPs[clientIP] = max(0, q.clientIPs[clientIP]-1)
	for i := slices.Index(q.hashes, hash); i != -1; i = slices.Index(q.hashes, hash) {
		q.hashes = slices.Delete(q.hashes, i, i+1)
	}
	q.mutex.Unlock()
}
