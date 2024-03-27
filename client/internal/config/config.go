package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Encryption EncryptionConfig
	Sync       SyncConfig
	Server     ServerConfig
}

type SyncConfig struct {
	Graphs   []string
	Interval int
	Once     bool
}
type EncryptionConfig struct {
	Enabled bool
	Key     string
}

type ServerConfig struct {
	Host     string
	ApiToken string
}

func Read() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("/etc/logsync")
	viper.AddConfigPath("$HOME/.logsync")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("LOGSYNC_CLIENT")
	viper.AutomaticEnv()

	defineDefaults()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	conf := getConfig()
	err = validateConfig(conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

func defineDefaults() {
	viper.SetDefault("encryption.enabled", false)
	viper.SetDefault("graphs", []string{})

	viper.SetDefault("sync.interval", 60)
	viper.SetDefault("sync.once", false)
}

func getConfig() Config {
	return Config{
		Encryption: EncryptionConfig{
			Enabled: viper.GetBool("encryption.enabled"),
			Key:     viper.GetString("encryption.key"),
		},
		Sync: SyncConfig{
			Graphs:   viper.GetStringSlice("sync.graphs"),
			Interval: viper.GetInt("sync.interval"),
			Once:     viper.GetBool("sync.once"),
		},
		Server: ServerConfig{
			Host:     viper.GetString("server.host"),
			ApiToken: viper.GetString("server.apitoken"),
		},
	}
}

func validateConfig(config Config) error {
	if config.Server.Host == "" {
		return errors.New("server.host is required")
	}

	if len(config.Sync.Graphs) == 0 {
		return errors.New("server.graphs must not be empty")
	}

	if config.Encryption.Enabled && config.Encryption.Key == "" {
		return errors.New("encryption.key is required, when encryption is enabled")
	}

	if !config.Sync.Once && config.Sync.Interval <= 0 {
		return errors.New("sync.interval must be set, when sync.once is disabled")
	}

	return nil
}
