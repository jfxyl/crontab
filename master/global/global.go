package global

import (
	"crontab/master/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ConfigPath string

	Config *config.Config

	EtcdClient *clientv3.Client

	MongoClient *mongo.Client

	LogCollection *mongo.Collection
)
