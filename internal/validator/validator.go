package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (v Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.addFieldError(key, message)
	}
}

func (v *Validator) AddFieldError(key, message string) {
	v.CheckField(false, key, message)
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Gte(value, floor int) bool {
	return value >= floor
}

func Lte(value, ceil int) bool {
	return value <= ceil
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
