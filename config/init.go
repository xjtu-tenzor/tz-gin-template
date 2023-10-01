package config

import (
	"template/service/validator"
)

func init() {
	initConfig()
	initLogger()
	validator.InitValidator()
}
