package jobstores

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go-Job-Scheduler/jobs"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	RedisKey    = "job::store"
	RuntimesKey = "job::runtimes"
)

type RedisJobStore struct {
	storeKey    string
	runtimesKey string
	Host        string
	Port        int
	DB          int
	password    string
	Client      *redis.Client
	sync.RWMutex
}

func newRedisJobStore() JobStore {
	return &RedisJobStore{
		storeKey:    RedisKey,
		runtimesKey: RuntimesKey,
	}
}

func (store *RedisJobStore) setOption(option StoreOption) {
	store.Host = option.Host
	port, _ := strconv.Atoi(option.Port)
	store.Port = port
	db, _ := strconv.Atoi(option.DBName)
	store.DB = db
	store.password = option.Password
	store.Client = redis.NewClient(&redis.Options{
		Addr:     store.Host + ":" + strconv.Itoa(store.Port),
		Password: store.password,
		DB:       store.DB,
	})
	err := store.Client.Ping().Err()
	if err != nil {
		panic(err)
	}
}

func (store *RedisJobStore) connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     store.Host + ":" + strconv.Itoa(store.Port),
		Password: store.password,
		DB:       store.DB,
	})
	defer func() {
		_ = client.Close()
	}()
}

func (store *RedisJobStore) AddJob(j jobs.Job) error {
	if store.Client.HExists(store.storeKey, j.Id).Val() {
		return errors.New(fmt.Sprint("job %s already exists", j.Id))
	}
	var job *jobs.Job
	// 如果传入的job id为空， 则调用jobs.New生成job id
	if strings.EqualFold(j.Id, "") {
		job = jobs.New(j.Name, j.FuncName, j.StartTime, j.Interval, j.Type, j.Args...)
	} else {
		job = &j
	}

	log.Println("ExecutionPeriodic=", jobs.ExecutionPeriodic)
	// 加锁
	store.Lock()
	// 函数执行完毕前解锁
	defer store.Unlock()
	// 准备pipeline，添加至redis中jobs hash表以及job激活时间的有序集合中
	pipe := store.Client.Pipeline()
	pipe.HSet(store.storeKey, job.Id, job.Bytes())
	// 如果job下次执行时间非0，则将该job下次执行时间作为Score放入redis有序集合中
	if job.NextRunTime() != 0 {
		pipe.ZAdd(store.runtimesKey, redis.Z{Score: job.NextRunTime(), Member: job.Id})
	}
	_, err := pipe.Exec()
	if err != nil {
		return errors.New(fmt.Sprintf("Error: RedisJobStore::AddJob, %s", err.Error()))
	}
	return nil
}

func (store *RedisJobStore) RemoveJob(job jobs.Job) error {
	// 加锁
	store.Lock()
	// 函数执行完毕前解锁
	defer store.Unlock()
	// 准备pipeline，将需删除的任务从redis中删除
	pipe := store.Client.Pipeline()
	pipe.HDel(store.storeKey, job.Id)
	pipe.ZRem(store.runtimesKey, job.Id)
	_, err := pipe.Exec()
	if err != nil {
		return errors.New(fmt.Sprintf("Error: RedisJobStore::RemoveJob, %s", err.Error()))
	}
	return nil
}

func (store *RedisJobStore) UpdateJob(job *jobs.Job, anotherJob jobs.Job) error {
	// 加锁
	store.Lock()
	// 函数执行完毕前解锁
	defer store.Unlock()
	err := job.Update(anotherJob)
	return err
}

func (store *RedisJobStore) GetJobById(id string) *jobs.Job {
	val := store.Client.HGet(store.storeKey, id).Val()
	return jobs.BytesToJob([]byte(val))
}

func (store *RedisJobStore) GetJobs2Run() []jobs.Job {
	var jobs2Run []jobs.Job
	now := time.Now().Unix()
	// 加锁
	store.Lock()
	// 函数执行完毕前解锁
	defer store.Unlock()

	// 从 redis 中按 Score 获取从0到当前时间戳之间的所有任务id
	results := store.Client.ZRangeByScore(store.runtimesKey, redis.ZRangeBy{
		Min: "1",
		Max: strconv.Itoa(int(now)),
	}).Val()

	if len(results) <= 0 {
		return jobs2Run
	}

	// 准备pipeline，将当前要执行的任务从redis中删除
	pipe := store.Client.Pipeline()
	for _, jobId := range results {
		job := store.GetJobById(jobId)
		if !strings.EqualFold(job.Id, "") {
			jobs2Run = append(jobs2Run, *job)
			pipe.HDel(store.storeKey, job.Id)
			pipe.ZRem(store.runtimesKey, job.Id)
		}
	}
	_, err := pipe.Exec()
	if err != nil {
		log.Println("Error: RedisJobStore::GetJobs2Run,", err)
	}
	return jobs2Run
}

func (store *RedisJobStore) GetAllJobs() []jobs.Job {
	var allJobs []jobs.Job
	results, err := store.Client.HGetAll(store.storeKey).Result()
	if err != nil {
		log.Println("Error: redisStore GetAllJobs, ", err)
	}
	for _, serializedStrJob := range results {
		job := jobs.BytesToJob([]byte(serializedStrJob))
		if job != nil {
			allJobs = append(allJobs, *job)
		}
	}
	return allJobs
}
