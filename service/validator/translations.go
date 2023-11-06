package validator

import (
	ut "github.com/go-playground/universal-translator"
)

func timingTransZh(ut ut.Translator) error {
	return ut.Add("timing", "{0}输入的时间不符合要求", true)
}
