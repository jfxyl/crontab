package controller

import (
	api "crontab/internal/api/service"
	"crontab/internal/common"
	"crontab/internal/job"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type JobController struct{}

func (c *JobController) Save(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		jobStr string
		newJob job.Job
		oldJob *job.Job
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	jobStr = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(jobStr), &newJob); err != nil {
		common.RespFail(w, err.Error())
		return
	}

	if oldJob, err = api.JobService.Save(&newJob); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, oldJob, "success")
}

func (c *JobController) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldJob *job.Job
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	name = r.PostForm.Get("name")
	if oldJob, err = api.JobService.Delete(name); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, oldJob, "success")
}

func (c *JobController) List(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		jobList []*job.Job
	)
	if jobList, err = api.JobService.List(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, jobList, "success")
}

func (c *JobController) Executings(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		name       string
		executings []string
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	name = r.Form.Get("name")
	if executings, err = api.JobService.Executings(name); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, executings, "success")
}

func (c *JobController) Kill(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		name      string
		timePoint string
		killJob   *job.KillJob
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	name = r.PostForm.Get("name")
	timePoint = r.PostForm.Get("time")
	killJob = &job.KillJob{
		Name: name,
		Time: timePoint,
	}
	if err = api.JobService.Kill(killJob); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, nil, "success")
}

func (c *JobController) Logs(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		name       string
		startParam string
		endParam   string
		pageParam  string
		limitParam string
		start      time.Time
		end        time.Time
		page       int
		limit      int
		filter     *common.JobLogFilter
		logs       []*job.JobLog
		total      int64
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	name = r.Form.Get("name")
	startParam = r.Form.Get("start")
	endParam = r.Form.Get("end")
	pageParam = r.Form.Get("page")
	limitParam = r.Form.Get("limit")
	if page, err = strconv.Atoi(pageParam); err != nil || page < 1 {
		page = 1
	}
	if limit, err = strconv.Atoi(limitParam); err != nil || limit < 1 {
		limit = 10
	}
	if start, err = time.ParseInLocation("2006-01-02 15:04:05", startParam, time.Local); err != nil {
		start = time.Time{}
	}
	if end, err = time.ParseInLocation("2006-01-02 15:04:05", endParam, time.Local); err != nil {
		end = time.Time{}
	}
	filter = &common.JobLogFilter{
		JobName: name,
		TimeRange: common.TimeRange{
			Start: start,
			End:   end,
		},
		Pagination: common.Pagination{
			Page:  page,
			Limit: limit,
		},
	}
	if logs, total, err = api.JobService.Logs(filter); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, map[string]any{
		"data":  logs,
		"total": total,
	}, "success")
}

func (c *JobController) WorkList(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		workList []string
	)
	if workList, err = api.JobService.WorkList(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, workList, "success")
}
