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

type validateHandle struct {
	validator.Func                     //校验规则
	validator.RegisterTranslationsFunc //翻译规则
}

// 自定义校验规则及翻译应在此处注册
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
			if err := zhTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
				panic(err)
			}
			for name, function := range validatorHandleRouter {
				err := v.RegisterValidation(name, function.Func)
				if err != nil {
					panic(err)
				}
			}
			for name, function := range validatorHandleRouter {
				if err := v.RegisterTranslation(name, Trans, function.RegisterTranslationsFunc, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T(name, fe.Field())
					return t
				}); err != nil {
					panic(err)
				}
			}
		case "en":
			validatorDefault(v)
		default:
			validatorDefault(v)
		}
	}
}

func validatorDefault(v *validator.Validate) {
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
