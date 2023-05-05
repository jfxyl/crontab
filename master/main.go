package main

import (
	"crontab/master/initialize"
	"log"
)

func main() {
	var (
		err error
	)
	//初始化配置文件路径
	initialize.InitArgs()
	//初始化配置文件
	if err = initialize.InitConfig(); err != nil {
		log.Fatalln(err)
	}
	//初始化etcd
	if err = initialize.InitEtcd(); err != nil {
		log.Fatalln(err)
	}
	//初始化mongodb
	if err = initialize.InitMongo(); err != nil {
		log.Fatalln(err)
	}
	//初始化http服务
	if err = initialize.InitHttpServer(); err != nil {
		log.Fatalln(err)
	}
}
