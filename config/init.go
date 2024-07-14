package config

import (
	"template/service/validator"
)

func Init() {
	initConfig()
	initLogger()
	validator.InitValidator(Config.AppLanguage)
}
