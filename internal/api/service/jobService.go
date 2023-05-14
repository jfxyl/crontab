package api

import (
	"context"
	"crontab/internal/common"
	"crontab/internal/global"
	"crontab/internal/job"
	"encoding/json"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

//func init() {
//	JobService = &jobService{}
//}

var (
	JobService *jobService
)

type jobService struct{}

func (s *jobService) Save(jobPtr *job.Job) (oldJob *job.Job, err error) {
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj job.Job
	)
	jobKey = common.JOB_SAVE_DIR + jobPtr.Name
	if jobValue, err = json.Marshal(jobPtr); err != nil {
		return
	}
	if putResp, err = global.EtcdClient.Put(context.Background(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

func (s *jobService) Delete(name string) (oldJob *job.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj job.Job
	)
	jobKey = common.JOB_SAVE_DIR + name
	if delResp, err = global.EtcdClient.Delete(context.Background(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(delResp.PrevKvs) > 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

func (s *jobService) List() (jobList []*job.Job, err error) {
	var (
		jobDir  string
		getResp *clientv3.GetResponse
		kv      *mvccpb.KeyValue
		jobPtr  *job.Job
	)
	jobDir = common.JOB_SAVE_DIR

	if getResp, err = global.EtcdClient.Get(context.Background(), jobDir, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kv = range getResp.Kvs {
		jobPtr = &job.Job{}
		if err = json.Unmarshal(kv.Value, jobPtr); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, jobPtr)
	}
	return
}

func (s *jobService) Executings(name string) (executings []string, err error) {
	var (
		lockingKey string
		getResp    *clientv3.GetResponse
		kv         *mvccpb.KeyValue
		timePoint  string
	)
	executings = make([]string, 0)
	lockingKey = common.JOB_LOCK_DIR + name

	if getResp, err = global.EtcdClient.Get(context.Background(), lockingKey, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kv = range getResp.Kvs {
		timePoint = strings.TrimLeft(string(kv.Key), lockingKey)
		executings = append(executings, timePoint)
	}
	return
}

func (s *jobService) Kill(killJob *job.KillJob) (err error) {
	var (
		killKey        string
		jobValue       []byte
		leaseGrantResp *clientv3.LeaseGrantResponse
	)
	if jobValue, err = json.Marshal(killJob); err != nil {
		return
	}
	killKey = common.JOB_KillER_DIR + killJob.Name
	if leaseGrantResp, err = global.EtcdClient.Grant(context.Background(), 10); err != nil {
		return
	}
	if _, err = global.EtcdClient.Put(context.Background(), killKey, string(jobValue), clientv3.WithLease(leaseGrantResp.ID)); err != nil {
		return
	}
	return
}

func (s *jobService) Logs(condition *common.JobLogFilter) (logs []*job.JobLog, total int64, err error) {
	var (
		opts   []*options.FindOptions
		filter bson.M
		cursor *mongo.Cursor
		jobLog *job.JobLog
	)
	logs = make([]*job.JobLog, 0)
	opts = make([]*options.FindOptions, 0)
	opts = append(opts, options.Find().SetSort(bson.D{{"startTime", -1}}))
	opts = append(opts, options.Find().SetSkip(int64((condition.Pagination.Page-1)*condition.Pagination.Limit)))
	opts = append(opts, options.Find().SetLimit(int64(condition.Pagination.Limit)))

	filter = bson.M{
		"jobName": condition.JobName,
	}
	if !condition.Start.IsZero() || !condition.End.IsZero() {
		timerange := bson.M{}
		if !condition.Start.IsZero() {
			timerange["$gte"] = condition.Start.Unix() * 1000
		}
		if !condition.End.IsZero() {
			timerange["$lte"] = condition.End.Unix()*1000 + 999
		}
		filter["startTime"] = timerange
	}
	if total, err = global.LogCollection.CountDocuments(context.Background(), filter); err != nil || total == 0 {
		return
	}
	if cursor, err = global.LogCollection.Find(context.Background(), filter, opts...); err != nil {
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		jobLog = &job.JobLog{}
		if err = cursor.Decode(jobLog); err != nil {
			continue
		}
		logs = append(logs, jobLog)
	}
	return
}

func (s *jobService) WorkList() (workList []string, err error) {
	var (
		getResp *clientv3.GetResponse
		kv      *mvccpb.KeyValue
		ip      string
	)
	if getResp, err = global.EtcdClient.Get(context.Background(), common.JOB_WORKER_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kv = range getResp.Kvs {
		ip = strings.TrimLeft(string(kv.Key), common.JOB_WORKER_DIR)
		workList = append(workList, ip)
	}
	return workList, nil
}
