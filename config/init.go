package config

import (
	"template/service/validator"
)

func init() {
	initConfig()
	initLogger()
	initRedis()
	validator.InitValidator(Config.AppLanguage)
}
