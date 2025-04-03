package config

import (
	"sync"



	"github.com/spf13/viper"
)

var appConfig AppConfig
var configReader *viper.Viper

// Conf 返回App配置基本对象
func Conf() *AppConfig {

	return &appConfig
}

// GetConfigReader 返回配置读取器
func ConfigReader() *viper.Viper {
	return configReader
}

// InitConfig 初始化配置,
//
// 应该在main函数中调用，且要在调用Conf()和GetConfigReader()之前执行,只执行一次。
func InitConfig(confPath,confName,confType string) {
	sync.OnceFunc(func() {
		v := viper.New()
		configReader = v
		configReader.SetConfigName(confName)
		configReader.AddConfigPath(confPath)
		configReader.SetConfigType(confType)
		err := configReader.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = configReader.Unmarshal(&appConfig)
	
		if err != nil {
			panic(err)
		}
	})
	
}


