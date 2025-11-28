package config

import "github.com/spf13/viper"

type Config struct {
	App        AppConfig
	DB         DBConfig
	Migrations MigrationConfig
}

type AppConfig struct {
	Port string
}

type DBConfig struct {
	URL        string
	SlavesUrls []string
}

type MigrationConfig struct {
	FilePath string
}

func NewConfig(confPath string) *Config {
	conf := viper.New()
	conf.SetConfigFile(confPath)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Config
	err = conf.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	config.DB.SlavesUrls = []string{config.DB.URL}
	config.Migrations.FilePath = conf.GetString("migrations.file_path")
	return &config
}
