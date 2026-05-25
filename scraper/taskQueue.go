package scraper

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/spf13/viper"
)

var enqueueScript = redis.NewScript(`
	local tasksKey, hashesKey, clientsKey = KEYS[1], KEYS[2], KEYS[3]
	local hash, taskJSON, taskClient, capacity, isFront = ARGV[1], ARGV[2], ARGV[3], tonumber(ARGV[4]), ARGV[5] == "1"

	if redis.call("SISMEMBER", hashesKey, hash) == 1 then return 0 end
	if not isFront and redis.call("LLEN", tasksKey) >= capacity then return -1 end

	redis.call("SADD", hashesKey, hash)
	if isFront then
		redis.call("LPUSH", tasksKey, taskJSON)
	else
		redis.call("RPUSH", tasksKey, taskJSON)
	end
	redis.call("HINCRBY", clientsKey, taskClient, 1)
	return 1
`)

type Task struct {
	Hash       string
	Metadata   map[string]string
	TaskClient string
	URL        string
}

type TaskQueue struct {
	capacity    int
	ctx         context.Context
	mutex       sync.Mutex
	processFunc func(Task)
	rdb         *redis.Client
	instanceID  string
	tasksKey    string
	hashesKey   string
	clientsKey  string
	pausedKey   string
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	opts, _ := redis.ParseURL(viper.GetString("redis"))

	// Explicitly disable maintenance notifications
	// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
	opts.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	rdb := redis.NewClient(opts)
	instanceID := viper.GetString("instanceid")

	queue := &TaskQueue{
		capacity:   bufferSize,
		ctx:        context.Background(),
		rdb:        rdb,
		instanceID: instanceID,
		tasksKey:   "scraper:" + instanceID + ":tasks",
		hashesKey:  "scraper:" + instanceID + ":hashes",
		clientsKey: "scraper:" + instanceID + ":clients",
		pausedKey:  "scraper:" + instanceID + ":paused",
	}

	// Unpause the queue on startup if it was left in a paused state
	rdb.Del(queue.ctx, queue.pausedKey)

	go queue.run()
	return queue
}

func (q *TaskQueue) AddTask(taskClient, hash, url string, front bool, metadata map[string]string) bool {
	task := Task{
		Hash:       hash,
		Metadata:   metadata,
		TaskClient: taskClient,
		URL:        url,
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return false
	}

	isFront := "0"
	if front {
		isFront = "1"
	}

	res, err := enqueueScript.Run(q.ctx, q.rdb, []string{q.tasksKey, q.hashesKey, q.clientsKey}, hash, taskJSON, taskClient, q.capacity, isFront).Int()
	return err == nil && res == 1
}

func (q *TaskQueue) run() {
	for {
		q.mutex.Lock()
		process := q.processFunc
		q.mutex.Unlock()

		// Ensure we don't pop tasks if there's no processor registered
		if process == nil {
			time.Sleep(time.Second)
			continue
		}

		// Check for pause state
		if paused, _ := q.rdb.Exists(q.ctx, q.pausedKey).Result(); paused > 0 {
			time.Sleep(time.Second)
			continue
		}

		// Blocking pop with timeout to allow periodic checks for the pause state
		res, err := q.rdb.BLPop(q.ctx, 5*time.Second, q.tasksKey).Result()
		if err != nil {
			continue
		}

		var task Task
		if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
			continue
		}

		process(task)
	}
}

func (q *TaskQueue) SetProcessFunc(f func(Task)) {
	q.mutex.Lock()
	q.processFunc = f
	q.mutex.Unlock()
}

func (q *TaskQueue) Pause(t time.Duration) {
	q.rdb.Set(q.ctx, q.pausedKey, "1", t)
}

func (q *TaskQueue) CountQueuedTasksForClient(taskClient string) (count int) {
	val, err := q.rdb.HGet(q.ctx, q.clientsKey, taskClient).Int()
	if err != nil {
		return 0
	}
	return max(0, val)
}

func (q *TaskQueue) ConfirmTaskCompletion(taskClient string, hash string) {
	pipe := q.rdb.Pipeline()
	pipe.HIncrBy(q.ctx, q.clientsKey, taskClient, -1)
	pipe.SRem(q.ctx, q.hashesKey, hash)
	_, _ = pipe.Exec(q.ctx)
}
