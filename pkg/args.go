package pkg

import (
	"flag"
	"fmt"

	"github.com/lbfatcgf/baseTemp/pkg/config"
)

var ConfigPath *string
var Port *string
var Version *bool

func ParseArgs() bool {
	configPath := flag.String("conf", "./conf", "config file directory")
	configName := flag.String("confname", "conf", "config file name")
	confType := flag.String("type", "yaml", "config file type")
	port := flag.String("port", "8888", "server port")
	vseions := flag.Bool("v", false, "show version")
	flag.Parse()
	ConfigPath = configPath
	Port = port
	Version = vseions
	config.InitConfig(*ConfigPath, *configName, *confType)
	if *Version {
		fmt.Println(config.Conf().Version)
		return true
	}
	return false
}
