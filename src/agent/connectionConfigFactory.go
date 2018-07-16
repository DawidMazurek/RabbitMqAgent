package main

import "github.com/spf13/viper"

func getConnectionConfig() *viper.Viper{
	config := viper.New()
	config.AutomaticEnv()
	config.SetEnvPrefix("rabbitmq")

	config.BindEnv("user")
	config.BindEnv("pass")
	config.BindEnv("host")
	config.BindEnv("port")

	return config
}