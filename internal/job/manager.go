package job

import (
	"context"
	"crontab/internal/common"
	"crontab/internal/global"
	"encoding/json"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	G_JobManager *Manager
)

func InitJobManager() {
	G_JobManager = &Manager{}
	go G_JobManager.WatchJob()
	go G_JobManager.WatchKiller()
}

//任务管理器
type Manager struct{}

func (s *Manager) WatchJob() (err error) {
	var (
		getResp    *clientv3.GetResponse
		kv         *mvccpb.KeyValue
		job        *Job
		jobEvent   *JobEvent
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
	)
	if getResp, err = global.EtcdClient.Get(context.Background(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kv = range getResp.Kvs {
		job = &Job{}
		if err = json.Unmarshal(kv.Value, job); err != nil {
			err = nil
			continue
		}
		jobEvent = G_JobScheduler.BuildJobEvent(common.JOB_EVENT_SAVE, job)
		G_JobScheduler.PushJobEvent(jobEvent)
	}
	go func() {
		watchChan = global.EtcdClient.Watch(context.Background(), common.JOB_SAVE_DIR, clientv3.WithPrefix(), clientv3.WithPrevKV())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if err = json.Unmarshal(watchEvent.Kv.Value, job); err != nil {
						continue
					}
					jobEvent = G_JobScheduler.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					if err = json.Unmarshal(watchEvent.PrevKv.Value, job); err != nil {
						continue
					}
					jobEvent = G_JobScheduler.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				G_JobScheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

func (s *Manager) WatchKiller() (err error) {
	var (
		job        *Job
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		killJob    *KillJob
		jobEvent   *JobEvent
	)
	go func() {
		watchChan = global.EtcdClient.Watch(context.Background(), common.JOB_KillER_DIR, clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					killJob = &KillJob{}
					if err = json.Unmarshal(watchEvent.Kv.Value, killJob); err != nil {
						continue
					}
					job = &Job{
						Name: killJob.Name,
					}
					jobEvent = G_JobScheduler.BuildKillJobEvent(common.JOB_EVENT_KILL, job, killJob.Time)
					G_JobScheduler.HandleJobEvent(jobEvent)
				case mvccpb.DELETE:
				}
			}
		}
	}()
	return
}
