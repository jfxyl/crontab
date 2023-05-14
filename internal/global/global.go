package global

import (
	"crontab/internal/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ConfigPath string

	Config *common.Config

	EtcdClient *clientv3.Client

	MongoClient *mongo.Client

	LogCollection *mongo.Collection
)
