package initialize

import (
	"crontab/master/global"
	"flag"
)

func InitArgs() {
	flag.StringVar(&global.ConfigPath, "config", "./master/config.yaml", "指定配置文件，默认./master/config.yaml")
	flag.Parse()
}
