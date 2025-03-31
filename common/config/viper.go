package config

import (
	"baseTemp/config"

	"github.com/spf13/viper"
)

var appConfig config.AppConfig

func Conf() *config.AppConfig {
	// fmt.Printf("config: %v\n", AppConfig)
	return &appConfig
}

func InitConfig(confPath string) {
	v := viper.New()
	v.SetConfigName("conf")
	v.AddConfigPath(confPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// for _, ak := range v.AllKeys() {
	// 	fmt.Println(v.Get(ak))
	// }

	err = v.Unmarshal(&appConfig)

	if err != nil {
		panic(err)
	}

	// fmt.Printf("config: %v\n", AppConfig)
}
