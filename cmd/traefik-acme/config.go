package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func configInit() {
	viper.SetConfigName("traefik")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./artifacts")
	viper.AddConfigPath("./test")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/traefik")
	viper.AddConfigPath("/usr/local/traefik/etc")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
