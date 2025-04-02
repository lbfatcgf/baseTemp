package cmd

import (
	"baseTemp/common/config"
	"flag"
	"fmt"
)

var ConfigPath *string
var Port *string
var Version *bool


func ParseArgs() bool {
	configPath := flag.String("conf", "./conf", "config file directory")
	port := flag.String("port", "8888", "server port")
	vseions := flag.Bool("v", false, "show version")
	flag.Parse()
	ConfigPath = configPath
	Port = port
	Version = vseions
	config.InitConfig(*ConfigPath)
	if(*Version){
		fmt.Println(config.Conf().Version)
		return true
	}
	return false
}