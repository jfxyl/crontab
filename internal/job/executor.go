package job

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

var (
	G_JobExecutor *Executor
)

func InitJobExecutor() {
	G_JobExecutor = &Executor{}
}

//任务执行器
type Executor struct{}

func (s *Executor) Run(jobExecuteInfo *JobExecuteInfo) {
	var (
		err    error
		cmd    *exec.Cmd
		output []byte
		result *JobExecuteResult = &JobExecuteResult{
			ExecuteInfo: jobExecuteInfo,
			Output:      make([]byte, 0),
		}
	)
	jobExecuteInfo.Lock = CreateJobLock(jobExecuteInfo)
	err = jobExecuteInfo.Lock.Lock()
	defer jobExecuteInfo.Lock.UnLock()
	result.StartTime = time.Now()
	if err != nil {
		result.Err = err
	} else {
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			cmd = exec.CommandContext(jobExecuteInfo.CancelCtx, "/bin/bash", "-c", jobExecuteInfo.Job.Command)
		} else if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(jobExecuteInfo.CancelCtx, "cmd", "/C", jobExecuteInfo.Job.Command)
		} else {
			log.Fatal("当前系统可能不支持")
		}
		output, err = cmd.CombinedOutput()
		fmt.Println("output", output)
		result.Output = output
		result.Err = err
	}
	result.EndTime = time.Now()
	G_JobScheduler.PushJobResult(result)
}

//func (s *Executor) PushExecuting(jobExecuteInfo *common.JobExecuteInfo) (err error) {
//	var (
//		executingKey string
//		//putResp      *clientv3.PutResponse
//	)
//	executingKey = fmt.Sprintf("%s%s/%s", common.JOB_Executing_DIR, jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime.Format("2006-01-02 15:04:05"))
//	if _, err = global.EtcdClient.Put(context.Background(), executingKey, ""); err != nil {
//		return
//	}
//	return
//}

type JobExecuteInfo struct {
	Job        *Job
	PlanTime   time.Time
	RealTime   time.Time
	Lock       *Lock
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

func (s *Executor) BuildJobExecuteInfo(jobSchedulerPlan *JobSchedulerPlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulerPlan.Job,
		PlanTime: jobSchedulerPlan.NextTime,
		RealTime: time.Now(),
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}
