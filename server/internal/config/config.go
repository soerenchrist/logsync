package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
)

type Config struct {
	Server  ServerConfig
	Files   FilesConfig
	Db      DbConfig
	Logging LoggingConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type FilesConfig struct {
	Path string
}

type DbConfig struct {
	Path string
}

type LoggingConfig struct {
	Level slog.Level
}

func (c Config) Url() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func Read() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("/etc/logsync")
	viper.AddConfigPath("$HOME/.logsync")
	viper.SetEnvPrefix("LOGSYNC")
	viper.SetEnvKeyReplacer(NewReplacer())
	viper.AutomaticEnv()

	defineDefaults()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	return getConfig(), nil
}

func defineDefaults() {
	viper.SetDefault("db.path", "logsync.db")
	viper.SetDefault("files.path", "files")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("log.level", "info")
}

func getConfig() Config {
	return Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Host: viper.GetString("server.host"),
		},
		Files: FilesConfig{
			Path: viper.GetString("files.path"),
		},
		Db: DbConfig{
			Path: viper.GetString("db.path"),
		},
		Logging: LoggingConfig{
			Level: getLogLevel(),
		},
	}
}

func getLogLevel() slog.Level {
	level := viper.GetString("log.level")
	switch level {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
