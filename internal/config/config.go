package config

import (
	"crontab/internal/common"
	"crontab/internal/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() (err error) {
	var (
		v *viper.Viper
	)
	global.Config = &common.Config{}
	v = viper.New()
	v.SetConfigFile(global.ConfigPath)
	//读取配置文件
	if err = readConfig(v); err != nil {
		return
	}
	//监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//编辑器可能会触发两次事件
		fmt.Printf("监听到文件变化：%s：%s", e.Name, e.Op)
		//读取配置文件
		if err = readConfig(v); err != nil {
			return
		}
		//do something
		fmt.Println("do something")
	})
	return
}

//读取配置文件
func readConfig(v *viper.Viper) (err error) {
	var (
		config common.Config
	)
	//读取配置文件
	if err = v.ReadInConfig(); err != nil {
		return
	}
	//解析配置文件(不直接使用global.Config去接收解析结果，在解析成功后再赋值，避免配置错误影响原来的配置)
	if err = v.Unmarshal(&config); err != nil {
		return
	}
	global.Config = &config
	return
}
