package config

import (
	"flag"
	"log"
	"os"
)

type HTTPServer struct {
	Addr string
}


type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

func MustLoad() {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "Path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
		log.Fatal("Config path is not set")
		}
		
	}
	
	
}