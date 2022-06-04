package api

import (
	"encoding/json"
	"go-Job-Scheduler/jobs"
	"go-Job-Scheduler/schedulers"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func jsonResponse(w http.ResponseWriter, resp *response) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

// route "/"，默认路由，404
func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", 404)
}

// route "/api/job/add"， 添加任务api
func handleJobAdd(w http.ResponseWriter, r *http.Request) {
	resp := &response{}
	defer func() {
		_ = jsonResponse(w, resp)
	}()

	err := r.ParseForm()
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}

	var j jobs.Job
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}

	scheduler := schedulers.GetScheduler()
	if !scheduler.IsRunning() {
		resp.Code = 1
		resp.Message = "scheduler is not running"
		return
	}

	scheduler.JobStore.AddJob(j)
	resp.Message = "success"
	return
}

// route "/api/jobs"，所有任务列表api
func handleJobsList(w http.ResponseWriter, r *http.Request) {
	resp := &response{}
	defer func() {
		_ = jsonResponse(w, resp)
	}()
	scheduler := schedulers.GetScheduler()
	if !scheduler.IsRunning() {
		resp.Code = 1
		resp.Message = "scheduler is not running"
		return
	}

	resp.Message = "success"
	resp.Data = scheduler.JobStore.GetAllJobs()
	return
}
