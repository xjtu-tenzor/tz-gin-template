package validator

import (
	ut "github.com/go-playground/universal-translator"
)

func categoryTransZh(ut ut.Translator) error {
	return ut.Add("category", "{0}输入格式或长度不符(十位二进制)", true)
}
