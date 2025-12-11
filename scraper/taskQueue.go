package scraper

import (
	"bdo-rest-api/utils"
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
	cond        *sync.Cond
	hashSet     map[string]struct{}
	mutex       sync.Mutex
	paused      bool
	processFunc func(Task)
	tasks       chan Task
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	queue := &TaskQueue{
		clientIPs: make(map[string]int),
		hashSet:   make(map[string]struct{}),
		tasks:     make(chan Task, bufferSize),
	}
	queue.cond = sync.NewCond(&queue.mutex)

	go queue.run()
	return queue
}

func (q *TaskQueue) AddTask(clientIP, hash, url string) (added bool) {
	fullURL := utils.BuildRequest(url, map[string]string{
		"taskClient": clientIP,
		"taskHash":   hash,
	})

	q.mutex.Lock()
	if _, exists := q.hashSet[hash]; exists {
		q.mutex.Unlock()
		return false
	}

	q.hashSet[hash] = struct{}{}
	q.clientIPs[clientIP]++
	q.mutex.Unlock()

	q.tasks <- Task{
		ClientIP: clientIP,
		Hash:     hash,
		URL:      fullURL,
	}
	return true
}

func (q *TaskQueue) run() {
	for task := range q.tasks {
		q.mutex.Lock()
		for q.paused {
			q.cond.Wait()
		}
		process := q.processFunc
		q.mutex.Unlock()

		if process != nil {
			process(task)
		}
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

	q.cond.Broadcast()
}

func (q *TaskQueue) CountQueuedTasksForClient(clientIP string) (count int) {
	q.mutex.Lock()
	count = max(0, q.clientIPs[clientIP])
	q.mutex.Unlock()

	return
}

func (q *TaskQueue) ConfirmTaskCompletion(clientIP string, hash string) {
	q.mutex.Lock()
	if q.clientIPs[clientIP] > 0 {
		q.clientIPs[clientIP]--
	}
	delete(q.hashSet, hash)
	q.mutex.Unlock()
}
