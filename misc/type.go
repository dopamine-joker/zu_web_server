package misc

type Config struct {
	Api     ApiConfig  `mapstructure:"api"`
	EtcdCfg EtcdConfig `mapstructure:"etcd"`
}

type ApiConfig struct {
	ListenPort int `mapstructure:"listenPort"`
}

type EtcdConfig struct {
	Host              []string `mapstructure:"host"`
	BasePath          string   `mapstructure:"basePath"`
	ServerPathLogic   string   `mapstructure:"serverPathLogic"`
	ServerPathConnect string   `mapstructure:"serverPathConnect"`
}
