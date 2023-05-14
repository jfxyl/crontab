package args

import (
	"crontab/internal/global"
	"flag"
)

func InitArgs() {
	flag.StringVar(&global.ConfigPath, "config", "./configs/config.yaml", "指定配置文件，默认../configs/config.yaml")
	flag.Parse()
}
