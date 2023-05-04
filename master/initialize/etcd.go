package initialize

import (
	"crontab/master/global"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func InitEtcd() (err error) {
	global.EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   global.Config.EtcdConfig.EndPoints,
		DialTimeout: time.Duration(global.Config.EtcdConfig.DialTimeout) * time.Second,
	})
	return
}
