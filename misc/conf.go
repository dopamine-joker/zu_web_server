package misc

import (
	"github.com/dopamine-joker/zu_web_server/db"
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
	viper.SetConfigName("config_local")
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
	db.InitRedis(Conf.RedisCfg.Address, Conf.RedisCfg.Port, Conf.RedisCfg.Password, Conf.RedisCfg.Db)
}

func initKey() {
	Key = os.Getenv("KEY")
}
