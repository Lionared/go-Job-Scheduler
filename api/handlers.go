package api

import (
	"encoding/json"
	"go-Job-Scheduler/jobs"
	"go-Job-Scheduler/schedulers"
	"net/http"
	"strings"
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

	err = scheduler.JobStore.AddJob(j)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}
	resp.Message = "success"
	return
}

// route "/api/job/delete, 删除任务api
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
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

	err = scheduler.JobStore.RemoveJob(j)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}
	resp.Message = "success"
	return
}

// route "/api/job/update, 修改任务api
func handleJobUpdate(w http.ResponseWriter, r *http.Request) {
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

	jobOld := scheduler.JobStore.GetJobById(j.Id)
	if jobOld == nil {
		resp.Code = 1
		resp.Message = "error: no such a job"
		return
	}

	err = scheduler.JobStore.UpdateJob(jobOld, j)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}
	resp.Message = "success"
	return
}

// route "/api/job/?id=xxx, 查询任务api
func handleJobRead(w http.ResponseWriter, r *http.Request) {
	resp := &response{}
	defer func() {
		_ = jsonResponse(w, resp)
	}()

	id := r.URL.Query().Get("id")
	if strings.EqualFold(id, "") {
		resp.Code = 1
		resp.Message = "must supply a job id param"
		return
	}

	scheduler := schedulers.GetScheduler()
	if !scheduler.IsRunning() {
		resp.Code = 1
		resp.Message = "scheduler is not running"
		return
	}

	job := scheduler.JobStore.GetJobById(id)
	if strings.EqualFold(job.Id, "") {
		resp.Code = 1
		resp.Message = "error: no such a job"
		return
	}
	resp.Message = "success"
	resp.Data = job
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
