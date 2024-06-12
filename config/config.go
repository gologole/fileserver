package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type ServerConfig struct {
	HTTP struct {
		Port              string
		ReadTimeout       time.Duration
		WriteTimeout      time.Duration
		ScheduledShutdown time.Duration
	}
	FTP struct {
		Port     string
		User     string
		Password string
	}
}

type DatabaseConfig struct {
	Driver string
	Source string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

var Conf Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Ошибка чтения конфигураций", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatal("Ошибка анмаршалинга конфигураций", err)
	}

	//бесполезное условие созданное для акцента внимания
	if Conf.Server.HTTP.ScheduledShutdown == 0 {
		Conf.Server.HTTP.ScheduledShutdown = 0
	}
}
