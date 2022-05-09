package config

import "github.com/spf13/viper"

type Config struct {
	ServerConfig serverConfig
}

type serverConfig struct {
	BasicAuthUser     string
	BasicAuthPassword string
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
