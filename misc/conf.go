package misc

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
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
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
	initLogger()
	initKey()
	initJaeger()
	fmt.Println(Conf)
}

func initKey() {
	Key = os.Getenv("KEY")
}
