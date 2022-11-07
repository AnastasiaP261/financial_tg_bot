package env

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
)

const (
	envLocal = "LOCAL"
	envProd  = "PROD"

	configFileLocal = "config/config_local.yaml"
	configFileProd  = "config/config_prod.yaml"
)

var inProd bool

// SetEnvVariable устанавливает значение переменной окружения - продовое или локальное.
// Должно быть вызвано в самом начале запуска приложения
func SetEnvVariable(env string) error {
	switch env {
	case envLocal:
		inProd = false
		logs.Info("app started in LOCAL env")
	case envProd:
		inProd = true
		logs.Info("app started in PROD env")
	default:
		return errors.New("invalid env variable")
	}

	return nil
}

func InProd() bool {
	return inProd
}

func GetConfigFilePath() (confPath string) {
	var configFile string
	if !inProd {
		configFile = configFileLocal
	} else {
		configFile = configFileProd
	}

	return configFile
}
