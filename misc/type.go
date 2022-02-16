package misc

type Config struct {
	Api       ApiConfig    `mapstructure:"api"`
	EtcdCfg   EtcdConfig   `mapstructure:"etcd"`
	JaegerCfg JaegerConfig `mapstructure:"jaeger"`
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

type JaegerConfig struct {
	Schema string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	Path   string `mapstructure:"path"`
}
