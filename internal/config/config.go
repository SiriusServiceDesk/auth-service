package config

import (
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type HttpServer struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"8000"`
}

type GrpcServer struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"8000"`
}

type Db struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Name     string `yaml:"db_name" env-default:"postgres"`
}

type Redis struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"6379"`
	Password string `yaml:"password"`
}

type JWT struct {
	Expires int    `yaml:"expires"`
	Secret  string `yaml:"secret"`
}

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"`
	GrpcServer GrpcServer `yaml:"grpc_server"`
	Db         Db         `yaml:"db"`
	Redis      Redis      `yaml:"redis"`
	Jwt        JWT        `yaml:"jwt"`
}

func GetConfig() *Config {

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		logger.Fatal("cannot read config file %s: %s", path, err)
	}

	return &config
}
