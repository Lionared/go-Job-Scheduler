package executors

import (
	"errors"
	"go-Job-Scheduler/jobs"
	"reflect"
)

const DefaultMaxPoolSize = 10

var (
	registeredFuncMap = make(map[string]interface{})
	executors         = make(map[string]Executor)
)

type Executor interface {
	Add(job jobs.Job)
	setOption(option ExecutorOption)
	Execute()
}

type ExecutorOption struct {
	PoolSize int
}

func call(funcName string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(registeredFuncMap[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("number of params is invalid")
		return
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func registerExecutors() {
	// 基础的执行器，将来可添加其他类型执行器
	executors["base"] = newBaseExecutor()
}

func init() {
	// 注册各种任务的执行函数
	registeredFuncMap["add"] = DoAdd
	registeredFuncMap["print"] = DoPrint
	// 注册各种执行器
	registerExecutors()
}

func MapToExecutorOption(m map[string]interface{}) ExecutorOption {
	var option ExecutorOption

	s := m["poolSize"]
	if v, ok := s.(int); ok {
		option.PoolSize = v
	} else {
		option.PoolSize = DefaultMaxPoolSize
	}
	return option
}

func NewExecutor(typeStr string, option ExecutorOption) Executor {
	if v, ok := executors[typeStr]; ok {
		v.setOption(option)
		return v
	}
	return nil
}
