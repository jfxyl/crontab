package job

import (
	"crontab/internal/common"
	"github.com/gorhill/cronexpr"
	"time"
)

var (
	G_JobScheduler *Scheduler
)

func InitJobScheduler() {
	G_JobScheduler = &Scheduler{
		JobEventChan:      make(chan *JobEvent, 1000),
		JobPlanTable:      make(map[string]*JobSchedulerPlan),
		JobExecutingTable: make(map[string]map[string]*JobExecuteInfo),
		JobResultChan:     make(chan *JobExecuteResult, 1000),
	}
	go G_JobScheduler.SchedulerLoop()
}

//任务调度器
type Scheduler struct {
	JobEventChan      chan *JobEvent
	JobPlanTable      map[string]*JobSchedulerPlan
	JobExecutingTable map[string]map[string]*JobExecuteInfo //map[任务名]map[时间点]*common.JobExecuteInfo
	JobResultChan     chan *JobExecuteResult
}

func (s *Scheduler) PushJobEvent(jobEvent *JobEvent) (err error) {
	s.JobEventChan <- jobEvent
	return
}

func (s *Scheduler) PushJobResult(jobResult *JobExecuteResult) (err error) {
	s.JobResultChan <- jobResult
	return
}

func (s *Scheduler) HandleJobEvent(jobEvent *JobEvent) (err error) {
	var (
		jobSchedulerPlan *JobSchedulerPlan
		jobExecuteInfo   *JobExecuteInfo
		jobExecuting     bool
		jobExisted       bool
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulerPlan, err = s.BuildJobSchedulerPlan(jobEvent.Job); err != nil {
			return
		}
		s.JobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulerPlan, jobExisted = s.JobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(s.JobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:
		if jobExecuteInfo, jobExecuting = s.JobExecutingTable[jobEvent.Job.Name][jobEvent.Time]; jobExecuting {
			jobExecuteInfo.CancelFunc()
			if jobExecuteInfo.Lock != nil {
				jobExecuteInfo.Lock.UnLock()
			}
		}
	}
	return
}

func (s *Scheduler) HandleJobResult(jobResult *JobExecuteResult) {
	var (
		jobLog *JobLog
	)
	delete(s.JobExecutingTable[jobResult.ExecuteInfo.Job.Name], jobResult.ExecuteInfo.PlanTime.Format("2006-01-02 15:04:05"))
	if jobResult.Err != common.ERR_LOCK_ALREADY_EXISTED {
		jobLog = &JobLog{
			JobName:       jobResult.ExecuteInfo.Job.Name,
			Command:       jobResult.ExecuteInfo.Job.Command,
			Output:        string(jobResult.Output),
			PlanTime:      jobResult.ExecuteInfo.PlanTime.UnixMilli(),
			SchedulerTime: jobResult.ExecuteInfo.RealTime.UnixMilli(),
			StartTime:     jobResult.StartTime.UnixMilli(),
			EndTime:       jobResult.EndTime.UnixMilli(),
		}
		if jobResult.Err != nil {
			jobLog.Err = jobResult.Err.Error()
		}
		G_JobLogger.PushJobLog(jobLog)
	}
}

func (s *Scheduler) SchedulerLoop() {
	var (
		jobEvent       *JobEvent
		jobResult      *JobExecuteResult
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
	)
	schedulerAfter = s.TryScheduler()
	schedulerTimer = time.NewTimer(schedulerAfter)
	for {
		select {
		case jobEvent = <-s.JobEventChan:
			s.HandleJobEvent(jobEvent)
			schedulerAfter = s.TryScheduler()
			schedulerTimer.Reset(schedulerAfter)
		case <-schedulerTimer.C:
			schedulerAfter = s.TryScheduler()
			schedulerTimer.Reset(schedulerAfter)
		case jobResult = <-s.JobResultChan:
			s.HandleJobResult(jobResult)
		}
	}
}

func (s *Scheduler) TryScheduler() (schedulerAfter time.Duration) {
	var (
		now              time.Time
		jobSchedulerPlan *JobSchedulerPlan
		jobExecuteInfo   *JobExecuteInfo
		nearTime         *time.Time
		timeExecuteMap   map[string]*JobExecuteInfo
		exists           bool
	)
	if len(s.JobPlanTable) == 0 {
		schedulerAfter = 1 * time.Second
		return
	}
	now = time.Now()
	for _, jobSchedulerPlan = range G_JobScheduler.JobPlanTable {
		if jobSchedulerPlan.NextTime.Equal(now) || jobSchedulerPlan.NextTime.Before(now) {
			jobExecuteInfo = G_JobExecutor.BuildJobExecuteInfo(jobSchedulerPlan)
			if timeExecuteMap, exists = s.JobExecutingTable[jobSchedulerPlan.Job.Name]; !exists {
				timeExecuteMap = make(map[string]*JobExecuteInfo)
				s.JobExecutingTable[jobSchedulerPlan.Job.Name] = timeExecuteMap
			}
			s.JobExecutingTable[jobSchedulerPlan.Job.Name][jobSchedulerPlan.NextTime.Format(common.YMDHIS)] = jobExecuteInfo
			go G_JobExecutor.Run(jobExecuteInfo)
			jobSchedulerPlan.NextTime = jobSchedulerPlan.Expr.Next(now)
		}
		if nearTime == nil || jobSchedulerPlan.NextTime.Before(*nearTime) {
			nearTime = &jobSchedulerPlan.NextTime
		}
	}
	schedulerAfter = (*nearTime).Sub(now)
	return
}

func (s *Scheduler) BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func (s *Scheduler) BuildKillJobEvent(eventType int, job *Job, time string) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Time:      time,
		Job:       job,
	}
}

func (s *Scheduler) BuildJobSchedulerPlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan, err error) {
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

type JobSchedulerPlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type JobEvent struct {
	EventType int
	Time      string
	Job       *Job
}

//指定kill某个时间点的任务
type KillJob struct {
	Name string `json:"name"`
	Time string `json:"time"`
}
