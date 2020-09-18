package client

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	Shortcut string `mapstructure:"Shortcut"`
}

func NewConfig() *Config {
	return setDefaults()
}

func setDefaults() *Config {
	var configuration Config

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("Shortcut", "SCROLL_LOCK")

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("failed unmarshalling default values: %s", err.Error())
	}

	return &configuration
}

func (c *Config) loadConfig() error {
	log.Info("loading config from config.toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Infof("failed to open config.toml, using defaults %s", err.Error())
		return err
	}

	if err := viper.UnmarshalExact(c); err != nil {
		log.Infof("unable to parse config.toml, using defaults: %s", err.Error())
		return err
	}

	log.Infof("using shortcut: %s", c.Shortcut)

	return nil
}
