package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
	LoggerConfig LoggerConfig
}

type ServerConfig struct {
	Host                     string
	Port                     string
	BasicAuthUser            string
	BasicAuthPassword        string
	CreateDashboardEndpoint  bool
	DeletePastesIntervalMins time.Duration
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type LoggerConfig struct {
	MinLevel    int
	Encoding    string
	PrintTraces bool
}

func LoadConfig(filepath string) (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigFile(filepath)
	//v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
