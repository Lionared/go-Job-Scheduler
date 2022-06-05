package executors

import (
	"go-Job-Scheduler/jobs"
	"log"
	"sync"
)

type BaseExecutor struct {
	PoolSize int
	Pool     []jobs.Job
}

func (this *BaseExecutor) setOption(option ExecutorOption) {
	this.PoolSize = option.PoolSize
}

func (this *BaseExecutor) Add(job jobs.Job) {
	this.Pool = append(this.Pool, job)
}

func (this *BaseExecutor) Execute() {
	var wg sync.WaitGroup
	var runningJobs int
	// TODO: 添加任务执行耗时、状态监测等
	for _, v := range this.Pool {
		if runningJobs < this.PoolSize {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Println("Executing job", v.Id)
				_, err := call(v.FuncName, v.Args...)
				if err != nil {
					log.Println("Error:", v.Id, err)
					// TODO: add failed job retries
				}
				log.Println("Executing job", v.Id, ". Done")
			}()
			runningJobs++
		}
	}
	wg.Wait()
	// 执行完清空 Pool
	this.Pool = this.Pool[0:0]
}

func newBaseExecutor() Executor {
	executor := &BaseExecutor{
		PoolSize: 10,
	}
	return executor
}
