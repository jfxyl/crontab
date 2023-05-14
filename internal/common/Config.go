package common

type etcdConfig struct {
	EndPoints   []string `mapstructure:"end_points" json:"end_points"`
	DialTimeout int64    `mapstructure:"dial_timeout" json:"dial_timeout"`
}

type mongoConfig struct {
	URI             string `mapstructure:"uri" json:"uri"`
	ConnnectTimeout string `mapstructure:"connnect_timeout" json:"connnect_timeout"`
}

type Config struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`

	EtcdConfig  etcdConfig  `mapstructure:"etcd" json:"etcd"`
	MongoConfig mongoConfig `mapstructure:"mongo" json:"mongo"`
}
