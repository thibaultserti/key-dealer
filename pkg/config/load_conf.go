package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"logLevel"`
	Hostname string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
}

func LoadConfig(config_file string) (config Configuration, err error) {
	viper := viper.New()

	viper.SetConfigFile(config_file)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("KEY_DEALER")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
