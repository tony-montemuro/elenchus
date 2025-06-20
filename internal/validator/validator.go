package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.addFieldError(key, message)
	}
}

func (v *Validator) AddError(key, message string) {
	v.CheckField(false, key, message)
}

func (v *Validator) addFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
