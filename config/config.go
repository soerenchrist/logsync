package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Encryption EncryptionConfig
	Graph      GraphConfig
}

type EncryptionConfig struct {
	enabled bool
	key     string
}

type GraphConfig struct {
	graphs []string
}

func Read() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/logsync")
	viper.AddConfigPath("$HOME/.logsync")

	defineDefaults()
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	return getConfig(), nil
}

func defineDefaults() {
	viper.SetDefault("encryption.enabled", false)
	viper.SetDefault("graphs", []string{})
}

func getConfig() Config {
	return Config{
		Encryption: EncryptionConfig{
			enabled: viper.GetBool("encryption.enabled"),
			key:     viper.GetString("encryption.key"),
		},
		Graph: GraphConfig{
			graphs: viper.GetStringSlice("graphs"),
		},
	}
}
