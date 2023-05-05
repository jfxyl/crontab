package service

import (
	"context"
	"crontab/common"
	"crontab/master/global"
	"encoding/json"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var JobService *jobService

func init() {
	JobService = &jobService{}
}

type jobService struct{}

func (s *jobService) Save(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)
	jobKey = common.JOB_SAVE_DIR + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
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

func (s *jobService) Delete(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
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

func (s *jobService) List() (jobList []*common.Job, err error) {
	var (
		jobDir  string
		getResp *clientv3.GetResponse
		kv      *mvccpb.KeyValue
		job     *common.Job
	)
	jobDir = common.JOB_SAVE_DIR

	if getResp, err = global.EtcdClient.Get(context.Background(), jobDir, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kv = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kv.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

func (s *jobService) Kill() {

}
