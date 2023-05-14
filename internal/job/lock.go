package job

import (
	"context"
	"crontab/internal/common"
	"crontab/internal/global"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//任务锁
type Lock struct {
	jobName    string
	time       time.Time
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	lockKey    string
	isLocked   bool
}

func CreateJobLock(jobExecuteInfo *JobExecuteInfo) *Lock {
	return &Lock{
		jobName:    jobExecuteInfo.Job.Name,
		time:       jobExecuteInfo.PlanTime,
		cancelFunc: nil,
	}
}

func (s *Lock) Lock() (err error) {
	var (
		ctx                    context.Context = context.TODO()
		cancelFunc             context.CancelFunc
		leaseGrantResp         *clientv3.LeaseGrantResponse
		leaseKeepAliveResp     *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveRespChan <-chan *clientv3.LeaseKeepAliveResponse
		txn                    clientv3.Txn
		txnResp                *clientv3.TxnResponse
		lockKey                string
	)
	ctx = context.TODO()
	//获取一个10s的租约
	if leaseGrantResp, err = global.EtcdClient.Grant(ctx, 10); err != nil {
		panic(err)
	}
	//自动续租
	ctx, cancelFunc = context.WithCancel(ctx)
	if leaseKeepAliveRespChan, err = global.EtcdClient.KeepAlive(ctx, leaseGrantResp.ID); err != nil {
		defer func() {
			cancelFunc()
			global.EtcdClient.Revoke(context.TODO(), leaseGrantResp.ID)
		}()
		return
	}
	go func() {
		for {
			select {
			case leaseKeepAliveResp = <-leaseKeepAliveRespChan:
				if leaseKeepAliveResp == nil {
					return
				}
			}
		}
	}()
	lockKey = fmt.Sprintf("%s%s/%s", common.JOB_LOCK_DIR, s.jobName, s.time.Format("2006-01-02 15:04:05"))
	txn = global.EtcdClient.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "1", clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(lockKey))
	if txnResp, err = txn.Commit(); err != nil {
		return
	}
	if !txnResp.Succeeded {
		err = common.ERR_LOCK_ALREADY_EXISTED
		return
	}
	s.cancelFunc = cancelFunc
	s.leaseId = leaseGrantResp.ID
	s.isLocked = true
	return
}

func (s *Lock) UnLock() {
	if s.isLocked {
		s.cancelFunc()
		global.EtcdClient.Revoke(context.TODO(), s.leaseId)
	}
}
