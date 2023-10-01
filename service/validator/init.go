package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validatorRouter = map[string]validator.Func{
	"timing": timing,
}

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for name, function := range validatorRouter {
			err := v.RegisterValidation(name, function)
			if err != nil {
				panic(err)
			}
		}
	}
}
