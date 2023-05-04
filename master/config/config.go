package config

type EtcdConfig struct {
	EndPoints   []string `mapstructure:"end_points" json:"end_points"`
	DialTimeout int64    `mapstructure:"dial_timeout" json:"dial_timeout"`
}

type MongoDBConfig struct {
	URI             string `mapstructure:"uri" json:"uri"`
	ConnnectTimeout string `mapstructure:"connnect_timeout" json:"connnect_timeout"`
}

type Config struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`

	EtcdConfig    EtcdConfig    `mapstructure:"etcd" json:"etcd"`
	MongoDBConfig MongoDBConfig `mapstructure:"mongodb" json:"mongodb"`
}
