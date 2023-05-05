package initialize

import (
	"crontab/common"
	"crontab/master/global"
	"crontab/master/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		jobStr string
		job    common.Job
		oldJob *common.Job
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	jobStr = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(jobStr), &job); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	if oldJob, err = service.JobService.Save(&job); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, oldJob, "success")
}

func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
	)
	if err = r.ParseForm(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	name = r.PostForm.Get("name")
	if oldJob, err = service.JobService.Delete(name); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, oldJob, "success")
}

func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		jobList []*common.Job
	)
	if jobList, err = service.JobService.List(); err != nil {
		common.RespFail(w, err.Error())
		return
	}
	common.RespOk(w, jobList, "success")
}

func InitHttpServer() (err error) {
	var (
		mux           *http.ServeMux
		staticHandler http.Handler
		httpServer    http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	//mux.HandleFunc("/job/kill", handleJobKill)
	//mux.HandleFunc("/job/log", handleJobLog)
	//mux.HandleFunc("/worker/list", handleWorkerList)

	staticHandler = http.FileServer(http.Dir("./master/resources"))
	mux.Handle("/", staticHandler)

	httpServer = http.Server{
		Addr:    fmt.Sprintf(":%d", global.Config.Port),
		Handler: mux,
	}
	err = httpServer.ListenAndServe()
	return
}
