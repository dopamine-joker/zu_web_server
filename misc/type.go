package misc

type Config struct {
	Api ApiConfig `mapstructure:"api"`
}

type ApiConfig struct {
	ListenPort int `mapstructure:"listenPort"`
}
