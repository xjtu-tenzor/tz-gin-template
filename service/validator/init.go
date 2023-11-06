package validator

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type validateHandle struct {
	validator.Func
	validator.RegisterTranslationsFunc
}

var validatorHandleRouter = map[string]validateHandle{
	"timing": {
		timing,
		timingTransZh,
	},
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
			validatorDefault(v, Trans)
			for name, function := range validatorHandleRouter {
				if err := v.RegisterTranslation(name, Trans, function.RegisterTranslationsFunc, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T(name, fe.Field())
					return t
				}); err != nil {
					panic(err)
				}
			}
		case "en":
			validatorDefault(v, Trans)
		default:
			validatorDefault(v, Trans)
		}
	}
}

func validatorDefault(v *validator.Validate, Trans ut.Translator) {
	if err := enTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
		panic(err)
	}
	for name, function := range validatorHandleRouter {
		err := v.RegisterValidation(name, function.Func)
		if err != nil {
			panic(err)
		}
	}
}
