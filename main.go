package main

import (
	"flag"
	"go-Job-Scheduler/api"
	"go-Job-Scheduler/schedulers"
	"log"
	"runtime"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {
	// 设置 goroutine 最大运行并发数
	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)
	// 解析启动参数
	var host string
	var port int
	var readTimeout int64
	var writeTimeout int64

	var storeType string
	var storeHost string
	var storePort string
	var storeDBName string
	var storeUsername string
	var storePassword string
	var storeCharset string

	var executorType string
	var executorPoolSize int
	flag.StringVar(&host, "h", "127.0.0.1", "-h, listening at 127.0.0.1 by default")
	flag.IntVar(&port, "p", 10028, "-p, listening at port 10027 by default")
	flag.Int64Var(&readTimeout, "rt", 5, "--rt, read timeout, default 5 seconds")
	flag.Int64Var(&writeTimeout, "wt", 60, "--wt, write timeout, default 60 seconds")

	flag.StringVar(&storeType, "store-type", "redis", "--store-type, job storage type, default is redis store")
	flag.StringVar(&storeHost, "store-host", "127.0.0.1", "--store-host")
	flag.StringVar(&storePort, "store-port", "0", "--store-port")
	flag.StringVar(&storeDBName, "store-dbname", "", "--store-dbname")
	flag.StringVar(&storeUsername, "store-username", "", "--store-username")
	flag.StringVar(&storePassword, "store-password", "", "--store-password")
	flag.StringVar(&storeCharset, "store-charset", "", "--store-charset")

	flag.StringVar(&executorType, "executor-type", "base", "--executor-type, job executor type, default is base executor")
	flag.IntVar(&executorPoolSize, "executor-pool-size", 10, "--executor-pool-size, default is 10")
	flag.Parse()

	// 初始化 scheduler
	scheduler := schedulers.NewScheduler(map[string]interface{}{
		"store": map[string]interface{}{
			"type": storeType,
			"options": map[string]interface{}{
				"host":     storeHost,
				"port":     storePort,
				"dbname":   storeDBName,
				"username": storeUsername,
				"password": storePassword,
				"charset":  storeCharset,
			},
		},
		"executor": map[string]interface{}{
			"type": executorType,
			"options": map[string]interface{}{
				"poolSize": executorPoolSize,
			},
		},
	})
	// 启动goroutine运行
	go scheduler.Run()
	// 启动web server
	server := api.NewWebServer(host, port, readTimeout, writeTimeout)
	server.Start()
}
