package config

import (
	"baseTemp/config"

	"github.com/spf13/viper"
)

var appConfig config.AppConfig
var configReader *viper.Viper

func Conf() *config.AppConfig {
	// fmt.Printf("config: %v\n", AppConfig)
	return &appConfig
}

func InitConfig(confPath string) {
	v := viper.New()
	configReader = v
	configReader.SetConfigName("conf")
	configReader.AddConfigPath(confPath)
	configReader.SetConfigType("yaml")
	err := configReader.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// for _, ak := range v.AllKeys() {
	// 	fmt.Println(v.Get(ak))
	// }

	err = configReader.Unmarshal(&appConfig)

	if err != nil {
		panic(err)
	}

	// fmt.Printf("config: %v\n", AppConfig)
}

func GetExpendConf(key string) any {
	return configReader.Get(key)
}
