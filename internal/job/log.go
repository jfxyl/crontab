package job

import (
	"context"
	"crontab/internal/global"
	"time"
)

//任务日志

var (
	G_JobLogger *Logger
)

type Logger struct {
	JobLogChan     chan *JobLog
	AutoCommitChan chan *LogBatch
}

func InitJobLogger() {
	G_JobLogger = &Logger{
		JobLogChan:     make(chan *JobLog),
		AutoCommitChan: make(chan *LogBatch),
	}
	go G_JobLogger.ProcessLoop()
}

func (s *Logger) PushJobLog(jobLog *JobLog) {
	s.JobLogChan <- jobLog
}

func (s *Logger) ProcessLoop() {
	var (
		jobLog       *JobLog
		logBatch     *LogBatch
		timeOutBatch *LogBatch
		commitTimer  *time.Timer
	)

	for {
		select {
		case jobLog = <-s.JobLogChan:
			if logBatch == nil {
				logBatch = &LogBatch{}
				commitTimer = time.AfterFunc(1*time.Second, func(logBatch *LogBatch) func() {
					return func() {
						s.AutoCommitChan <- logBatch
					}
				}(logBatch))
			}
			logBatch.Logs = append(logBatch.Logs, jobLog)
			if len(logBatch.Logs) >= 100 {
				s.save(logBatch)
				logBatch = nil
				commitTimer.Stop()
			}
		case timeOutBatch = <-s.AutoCommitChan:
			if timeOutBatch != logBatch {
				continue
			}
			s.save(logBatch)
			logBatch = nil
		}
	}
}

func (s *Logger) save(logBatch *LogBatch) {
	global.LogCollection.InsertMany(context.Background(), logBatch.Logs)
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
	Logs []any
}
