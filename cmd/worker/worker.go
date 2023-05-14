package main

import (
	"crontab/internal/args"
	"crontab/internal/config"
	"crontab/internal/etcd"
	"crontab/internal/job"
	"crontab/internal/mongo"
	"log"
	"time"
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
	//服务注册
	if err = job.InitRegister(); err != nil {
		log.Fatalln(err)
	}
	//初始化任务Logger
	job.InitJobLogger()
	//初始化任务执行器
	job.InitJobExecutor()
	//初始化任务调度器
	job.InitJobScheduler()
	//初始化任务管理器
	job.InitJobManager()
	for {
		time.Sleep(1 * time.Second)
	}
}
