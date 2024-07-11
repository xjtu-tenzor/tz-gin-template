package validator

import (
	ut "github.com/go-playground/universal-translator"
)

func timingTransZh(ut ut.Translator) error {
	return ut.Add("timing", "{0}输入的时间不符合要求", true)
}

func phoneTransZh(ut ut.Translator) error {
	return ut.Add("phone", "{0}输入的手机号码不符合要求", true)
}

func emailTransZh(ut ut.Translator) error {
	return ut.Add("email", "{0}输入的邮箱地址不符合要求", true)
}
