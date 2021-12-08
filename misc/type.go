package misc

type Config struct {
	Api ApiConfig `mapstruct:"api"`
}

type ApiConfig struct {
	ListenPort int `mapstruct:"listenPort"`
}
