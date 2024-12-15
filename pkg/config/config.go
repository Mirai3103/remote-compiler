package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Executor ExecutorConfig
	NATS     NATSConfig
	GRPC     GRPCConfig
}

type ServerConfig struct {
	Port int
}

type ExecutorConfig struct {
	IsolateCommand string
	IsolateBoxPath string
	IsolateDir     string
}

type NATSConfig struct {
	URL string
}

type GRPCConfig struct {
	Port int
}

func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("executor.isolateCommand", "isolate")
	viper.SetDefault("executor.isolateBoxPath", "/var/local/lib/isolate")
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("grpc.port", 50051)
	viper.SetDefault("executor.isolateDir", "/isolateBox")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
