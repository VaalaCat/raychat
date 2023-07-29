package settings

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type RayConfig struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	Email        string `env:"EMAIL"`
	Password     string `env:"PASSWORD"`
	Token        string `env:"TOKEN" env-default:""`
}

var rayConf RayConfig

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Warn("load .env file error, try to read from env")
	}
	err := cleanenv.ReadEnv(&rayConf)
	if err != nil {
		logrus.Panic("read env error", err)
	}
}

func Get() RayConfig {
	return rayConf
}
