package main

import (
	"crontab/internal/api"
	"crontab/internal/args"
	"crontab/internal/config"
	"crontab/internal/etcd"
	"crontab/internal/mongo"
	"log"
)

func main() {
	var (
		err error
	)
	//初始化配置文件路径
	args.InitArgs()
	//初始化配置文件
	if err = config.InitConfig(); err != nil {
		log.Fatalln(err)
	}
	//初始化etcd
	if err = etcd.InitEtcd(); err != nil {
		log.Fatalln(err)
	}
	//初始化mongodb
	if err = mongo.InitMongo(); err != nil {
		log.Fatalln(err)
	}
	//初始化http服务
	if err = api.InitHttpServer(); err != nil {
		log.Fatalln(err)
	}
}
