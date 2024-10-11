package config

import (
	"github.com/spf13/viper"
	"log"
)

var LocalConf *Config

type (
	Config struct {
		Engine Engine
	}
	Engine struct {
		Http Http
	}
	Http struct {
		Addr string
	}
)

func init() {
	LocalConf = &Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = viper.Unmarshal(LocalConf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
