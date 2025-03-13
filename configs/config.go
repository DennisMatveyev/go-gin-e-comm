package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env    string
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

func MustLoadConfig(path string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	return config
}
