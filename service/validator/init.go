package validator

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var validatorRouter = map[string]validator.Func{
	"timing": timing,
}

var validatorTransRouter = map[string]validator.RegisterTranslationsFunc{
	"category": categoryTransZh,
}

var Trans ut.Translator

func InitValidator(locale string) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		enT := en.New()

		uni := ut.New(enT, zhT, enT)

		var ok bool
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			panic(fmt.Errorf("uni.GetTranslator(%s) failed", locale))
		}

		switch locale {
		case "zh":
			if err := zhTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
				panic(err)
			}
			for name, function := range validatorRouter {
				err := v.RegisterValidation(name, function)
				if err != nil {
					panic(err)
				}
			}
			for name, function := range validatorTransRouter {
				if err := v.RegisterTranslation(name, Trans, function, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T(name, fe.Field())
					return t
				}); err != nil {
					panic(err)
				}
			}
		case "en":
			if err := enTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
				panic(err)
			}
			for name, function := range validatorRouter {
				err := v.RegisterValidation(name, function)
				if err != nil {
					panic(err)
				}
			}
		default:
			if err := enTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
				panic(err)
			}
			for name, function := range validatorRouter {
				err := v.RegisterValidation(name, function)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
