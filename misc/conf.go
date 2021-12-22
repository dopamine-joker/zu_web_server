package misc

import (
	"github.com/spf13/viper"
)

var (
	Conf Config
	Key  string
)

func Init() {
	var err error
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	initLogger()
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
