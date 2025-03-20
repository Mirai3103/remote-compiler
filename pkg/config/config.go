package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Executor ExecutorConfig `json:"executor"`
	GRPC     GRPCConfig     `json:"grpc"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int `json:"port"`
}

// ExecutorConfig holds executor-related configuration
type ExecutorConfig struct {
	IsolateDir           string `json:"isolateDir"`
	Strategy             string `json:"strategy"`
	MaxCompileConcurrent int    `json:"max_compile_concurrent"`
	MaxExecuteConcurrent int    `json:"max_execute_concurrent"`
}

// GRPCConfig holds gRPC-related configuration
type GRPCConfig struct {
	Port int `json:"port"`
}

// Default configuration values
var defaultConfig = Config{
	Server: ServerConfig{
		Port: 8080,
	},
	Executor: ExecutorConfig{
		IsolateDir: "./temp",
		Strategy:   "risk",
	},
	GRPC: GRPCConfig{
		Port: 50051,
	},
}

// LoadConfig loads configuration from file, creating it if it doesn't exist
func LoadConfig(filename string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("json")

	// Set default values
	v.SetDefault("server.port", defaultConfig.Server.Port)
	v.SetDefault("grpc.port", defaultConfig.GRPC.Port)
	v.SetDefault("executor.strategy", defaultConfig.Executor.Strategy)
	v.SetDefault("executor.isolateDir", defaultConfig.Executor.IsolateDir)
	v.SetDefault("executor.maxCompileConcurrent", 10)
	v.SetDefault("executor.maxExecuteConcurrent", 10)
	//viper.SetDefault("executor.isolateBoxPath", "/var/local/lib/isolate")

	// Try to read the config file
	err := v.ReadInConfig()
	if err != nil {
		// If file not found, create it with default values
		var configFileNotFoundError viper.ConfigFileNotFoundError
		var pathError *fs.PathError
		if errors.As(err, &configFileNotFoundError) || errors.As(err, &pathError) {
			// Ensure directory exists
			dir := filepath.Dir(filename)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}

			// Create config with default values
			config := defaultConfig

			// Create the file
			v.AddConfigPath(dir)
			if err := v.SafeWriteConfig(); err != nil {
				// If file already exists, just write to it
				if os.IsExist(err) {
					if err := v.WriteConfig(); err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			return &config, nil
		}

		return nil, err
	}

	// Unmarshal the config into our struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
