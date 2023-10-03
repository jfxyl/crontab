package args

import (
	"crontab/internal/global"
	"flag"
)

func InitArgs() {
	flag.StringVar(&global.ConfigPath, "config", "./config.yaml", "指定配置文件，默认../config.yaml")
	flag.Parse()
}
