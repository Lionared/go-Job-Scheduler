package schedulers

import (
	"go-Job-Scheduler/executors"
	"go-Job-Scheduler/jobs"
	"go-Job-Scheduler/jobstores"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var scheduler *baseScheduler

type baseScheduler struct {
	running  bool
	JobStore jobstores.JobStore
	Executor executors.Executor
}

func NewScheduler(m map[string]interface{}) *baseScheduler {
	var jobStore jobstores.JobStore
	var executor executors.Executor

	s := GetScheduler()

	// 配置 job store
	storeConfig := m["store"]
	if v, ok := storeConfig.(map[string]interface{}); ok {
		storeType := v["type"].(string)
		optionsMap := v["options"].(map[string]interface{})
		options := jobstores.MapToStoreOption(optionsMap)
		jobStore = jobstores.NewJobStore(storeType, options)
	} else {
		return nil
	}

	// 配置 job executor
	executorConfig := m["executor"]
	if v, ok := executorConfig.(map[string]interface{}); ok {
		executorType := v["type"].(string)
		optionsMap := v["options"].(map[string]interface{})
		options := executors.MapToExecutorOption(optionsMap)
		executor = executors.NewExecutor(executorType, options)
	} else {
		return nil
	}
	s.JobStore = jobStore
	s.running = true
	s.Executor = executor

	return s
}

func GetScheduler() *baseScheduler {
	return scheduler
}

func (this *baseScheduler) IsRunning() bool {
	return this.running
}

func (this *baseScheduler) Run() {
	// 以秒为单位的ticker
	var ticker = time.NewTicker(time.Second * 1)
	// 监听程序结束信号
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Interrupt,
		os.Kill,
	)
	for {
		select {
		case <-ticker.C:
			// 调度
			if this.running == true {
				// 获取到期应执行任务
				jobs2Run := this.JobStore.GetJobs2Run()
				// 遍历
				for _, job := range jobs2Run {
					// 将任务交给executor
					this.Executor.Add(job)
					// 如果为周期性任务
					if job.Type == jobs.ExecutionPeriodic {
						// 修改任务的下次执行时间
						job.NextRunTime_ = job.NextRunTime_.Add(job.Interval)
					} else {
						// 一次性任务则将下次执行时间设置为0
						job.NextRunTime_ = time.Unix(0, 0)
					}
					// 将任务放回 store
					this.JobStore.AddJob(job)
				}
				// 另起一个 goroutine 执行 executor
				go this.Executor.Execute()
			}
		case <-signalChan:
			this.running = false
			ticker.Stop()
			log.Println("scheduler stop")
			break
		}
	}
}

func init() {
	// 调度器设置为单例模式
	var once sync.Once
	once.Do(func() {
		scheduler = &baseScheduler{}
	})
}
