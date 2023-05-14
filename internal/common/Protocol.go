package common

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type YmdHisTime time.Time

func (t YmdHisTime) MarshalJSON() ([]byte, error) {
	n := time.Time(t)
	if n.IsZero() {
		return []byte("null"), nil
	}
	return []byte(n.Format(YMDHIS)), nil
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type JobLogFilter struct {
	JobName string `bson:"jobName"`
	TimeRange
	Pagination
}

type TimeRange struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}

type Pagination struct {
	Page  int
	Limit int
}

type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
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
