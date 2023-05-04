package global

import (
	"crontab/master/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	ConfigPath string

	Config *config.Config

	EtcdClient *clientv3.Client
)
