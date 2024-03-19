package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Encryption EncryptionConfig
	Sync       SyncConfig
}

type SyncConfig struct {
	Graphs   []string
	Interval int
}
type EncryptionConfig struct {
	enabled bool
	key     string
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

	viper.SetDefault("sync.interval", 60)
}

func getConfig() Config {
	return Config{
		Encryption: EncryptionConfig{
			enabled: viper.GetBool("encryption.enabled"),
			key:     viper.GetString("encryption.key"),
		},
		Sync: SyncConfig{
			Graphs:   viper.GetStringSlice("sync.graphs"),
			Interval: viper.GetInt("sync.interval"),
		},
	}
}
