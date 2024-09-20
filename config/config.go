package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	errFailedToLoadConfig = "error on parsing configuration file"
	errNotInited          = "config is not inited"
)

type consul struct {
	Enable      bool
	ServiceName string
	Tags        []string
}

type tracer struct {
	Name     string
	Endpoint string
}

type conf struct {
	Consul      consul
	Environment string `json:"environment"`
	Tracer      tracer
}

// Config ...
var Config conf

var config *viper.Viper

func init() {
	println("loading env...")

	Load()
}

// Load is for load config by env
func Load() {
	println("initing config...")

	var err error

	config = viper.New()
	// config.SetConfigType("yaml")
    config.SetConfigType("yaml")
    config.SetConfigName("config")
    config.AddConfigPath(".")
    config.AddConfigPath("../")
    config.AddConfigPath("../../")
	config.AddConfigPath("../../../")
	config.AddConfigPath("../../../../")

	err = config.ReadInConfig()
	if err != nil {
		log.Fatalf("%s - <%s>\n", errFailedToLoadConfig, err)
	}
	err = config.Unmarshal(&Config)
	if err != nil {
		err = fmt.Errorf("配置文件序列化失败: %v", err)
		panic(err)
	}

	return
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

// GetConfig get config
func GetConfig() *viper.Viper {
	if config == nil {
		log.Fatalln(errNotInited)
	}

	return config
}
