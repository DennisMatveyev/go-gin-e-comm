package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
		Mode string
	}
	Database struct {
		Uri string
	}
	JWT struct {
		Secret string
	}
	Log struct {
		Path       string
		MaxSize    int
		MaxBackups int
		MaxAge     int
		Compress   bool
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
