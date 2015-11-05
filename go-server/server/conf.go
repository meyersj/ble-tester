package server

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

const DEFAULT_PORT string = "8082"
const DEFAULT_REDIS string = "127.0.0.1:6379"

type Config struct {
	Port  string
	Redis string
}

// create and return Config data
func Read_config(filename string) *Config {
	var conf Config
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		conf.Port = DEFAULT_PORT
		conf.Redis = DEFAULT_REDIS
		fmt.Println("Failure processing config file: ", filename)
		fmt.Println("  default Port =", conf.Port)
		fmt.Println("  default Redis =", conf.Redis)
	}
	return &conf
}
