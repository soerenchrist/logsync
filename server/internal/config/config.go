package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Files  FilesConfig
	Db     DbConfig
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

	defineDefaults()
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	return getConfig(), nil
}

func defineDefaults() {
	viper.SetDefault("db.path", "logsync.db")
	viper.SetDefault("files.path", "files")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 3000)
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
	}
}
