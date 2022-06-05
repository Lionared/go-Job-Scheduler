package jobstores

import (
	"go-Job-Scheduler/jobs"
	"strconv"
	"sync"
)

type JobStore interface {
	setOption(StoreOption)
	AddJob(jobs.Job) error
	RemoveJob(jobs.Job) error
	UpdateJob(*jobs.Job, jobs.Job) error
	GetJobById(string) *jobs.Job
	GetJobs2Run() []jobs.Job
	GetAllJobs() []jobs.Job
	sync.Locker
}

type StoreOption struct {
	Host     string
	Port     string
	DBName   string
	CharSet  string
	Username string
	Password string
}

var jobStores = make(map[string]JobStore)

func init() {
	registerStores()
}

func MapToStoreOption(m map[string]interface{}) StoreOption {
	var options StoreOption

	if option, exist := m["host"]; exist {
		if v, ok := option.(string); ok {
			options.Host = v
		} else {
			options.Host = "127.0.0.1"
		}
	}

	if option, exist := m["port"]; exist {
		if v, ok := option.(string); ok {
			options.Port = v
		} else {
			if i, ok := option.(int); ok {
				options.Port = strconv.Itoa(i)
			} else {
				options.Port = "0"
			}
		}
	}

	if option, exist := m["dbname"]; exist {
		if v, ok := option.(string); ok {
			options.DBName = v
		} else {
			if i, ok := option.(int); ok {
				options.DBName = strconv.Itoa(i)
			}
		}
	}

	if option, exist := m["username"]; exist {
		if v, ok := option.(string); ok {
			options.Username = v
		} else {
			options.Username = ""
		}
	}

	if option, exist := m["password"]; exist {
		if v, ok := option.(string); ok {
			options.Password = v
		} else {
			options.Password = ""
		}
	}

	if option, exist := m["charset"]; exist {
		if v, ok := option.(string); ok {
			options.CharSet = v
		} else {
			options.CharSet = ""
		}
	}
	return options
}

func registerStores() {
	// 注册各种store，将来可添加 sql store, zookeeper store等
	jobStores["redis"] = newRedisJobStore()
}

func NewJobStore(typeStr string, option StoreOption) JobStore {
	if v, ok := jobStores[typeStr]; ok {
		v.setOption(option)
		return v
	}
	return nil
}
