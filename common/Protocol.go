package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"log"
	"net/http"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type JobEvent struct {
	EventType int
	Job       *Job
}

type JobSchedulerPlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type JobExecuteInfo struct {
	Job        *Job
	PlanTime   time.Time
	RealTime   time.Time
	CancelCtx  context.Context
	CancelFunc context.CancelFunc
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
}

type JobLog struct {
	JobName       string `bson:"jobName" json:"jobName"`
	Command       string `bson:"command" json:"command"`
	Err           string `bson:"err" json:"err"`
	Output        string `bson:"output" json:"output"`
	PlanTime      int64  `bson:"planTime" json:"planTime"`
	SchedulerTime int64  `bson:"schedulerTime" json:"schedulerTime"`
	StartTime     int64  `bson:"startTime" json:"startTime"`
	EndTime       int64  `bson:"endTime" json:"endTime"`
}

type LogBatch struct {
	Logs []interface{}
}

type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data
	resp, err = json.Marshal(response)
	return
}

func Resp(w http.ResponseWriter, errno int, data interface{}, msg string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)

	resp := Response{
		Errno: errno,
		Msg:   msg,
		Data:  data,
	}
	res, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("json 失败！")
	} else {
		w.Write(res)
	}
}

func RespOk(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, 1, nil, msg)
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func ExtractKillerName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_KillER_DIR)
}

func ExtractWorkerIp(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulerPlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	jobSchedulerPlan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

func BuildJobExecuteInfo(jobSchedulerPlan *JobSchedulerPlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulerPlan.Job,
		PlanTime: jobSchedulerPlan.NextTime,
		RealTime: time.Now(),
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}
