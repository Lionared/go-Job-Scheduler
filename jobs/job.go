package jobs

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"time"
)

const (
	ExecutionOnce = 1 << iota
	ExecutionPeriodic
)

const (
	ParseTimeLayout = "2006-01-02 15:04:05"
)

type Job struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	FuncName     string        `json:"funcName"`
	Args         []interface{} `json:"args"`
	StartTime    time.Time     `json:"startTime"`
	NextRunTime_ time.Time
	Interval     time.Duration `json:"interval"`
	Type         uint8         `json:"type"`
}

// New returns a valid job
// @param name : job name,
// @param funcName: 要执行的函数名称
// @param startTime: 任务开始时间, web传递格式为 2022-06-03T18:02:03Z
// @param interval: 周期性任务执行时间间隔，秒
// @param jobType: 1 一次性任务，2 周期性任务
// @param args: 要执行函数的参数
func New(name, funcName string, startTime time.Time, interval time.Duration, jobType uint8, args ...interface{}) *Job {
	id := uuid.New().String()
	_startTime, _ := time.ParseInLocation(ParseTimeLayout, startTime.Format(ParseTimeLayout), time.Local)
	return &Job{
		Id:           id,
		Name:         name,
		FuncName:     funcName,
		Args:         args,
		StartTime:    _startTime,
		NextRunTime_: _startTime,
		Interval:     interval * time.Second,
		Type:         jobType,
	}
}

func (job *Job) NextRunTime() float64 {
	t := job.NextRunTime_.Unix()
	return float64(t)
}

func (job *Job) Update(modified Job) error {
	if modified.Name != "" {
		job.Name = modified.Name
	}
	if modified.Interval != 0 {
		job.Interval = modified.Interval
	}
	return nil
}

func (job *Job) String() string {
	return "Job:" + job.Name + ":" + job.Id
}

func (job *Job) Bytes() []byte {
	// 使用 encoding/gob 序列化
	buf := new(bytes.Buffer)
	_ = gob.NewEncoder(buf).Encode(job)
	return buf.Bytes()
}

func BytesToJob(b []byte) *Job {
	// 使用 encoding/gob 反序列化
	var job Job
	_ = gob.NewDecoder(bytes.NewBuffer(b)).Decode(&job)
	return &job
}
