package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DBPath  string `mapstructure:"db_path"`
	DataDir string `mapstructure:"data_dir"`
}

var AppConfig Config

func Load() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dataDir := filepath.Join(home, ".local", "share", "tkt")
	os.MkdirAll(dataDir, 0755)

	viper.SetDefault("data_dir", dataDir)
	viper.SetDefault("db_path", filepath.Join(dataDir, "tkt.db"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dataDir)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// Config file not found → use defaults
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}

func GetDBPath() string {
	return AppConfig.DBPath
}

func GetDataDir() string {
	return AppConfig.DataDir
}
