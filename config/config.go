package config

import "github.com/spf13/viper"

type Config struct {
	ServerConfig serverConfig
	DBConfig     dbConfig
}

type serverConfig struct {
	Host                    string
	Port                    string
	BasicAuthUser           string
	BasicAuthPassword       string
	CreateDashboardEndpoint bool
}

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
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
