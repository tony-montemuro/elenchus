package validator

import (
	"fmt"
)

type RangeRule struct {
	MinLength int
	MaxLength int
}

type RangeError struct {
	Key     string
	message string
}

func (e *RangeError) Error() string {
	return e.message
}

type FormRangeRules map[string]RangeRule

var RangeRules = map[string]FormRangeRules{
	SignUpForm: {
		"firstName": {MinLength: 1, MaxLength: 50},
		"lastName":  {MinLength: 1, MaxLength: 50},
		"email":     {MinLength: 1, MaxLength: 255},
		"password":  {MinLength: 8, MaxLength: 72},
	},
	LoginForm: {
		"email":    {MinLength: 1, MaxLength: 255},
		"password": {MinLength: 8, MaxLength: 72},
	},
}

func GetRangeErrors(form Form, name string) []RangeError {
	rules := RangeRules[name]
	errs := []RangeError{}

	for key, val := range form.GetStringVals() {
		var err *RangeError
		rule, exists := rules[key]
		if exists {
			err = getError(rule, key, val)
		}

		if err != nil {
			errs = append(errs, *err)
		}
	}

	return errs
}

func getError(rule RangeRule, formKey, formValue string) *RangeError {
	mn, mx := rule.MinLength, rule.MaxLength

	if mn > 0 {
		if mn == 1 && !NotBlank(formValue) {
			return &RangeError{Key: formKey, message: "This field cannot be blank."}
		} else if !MinChars(formValue, mn) {
			return &RangeError{Key: formKey, message: fmt.Sprintf("This field must be at least %d characters long.", mn)}
		}
	}
	if mx > 0 && !MaxChars(formValue, mx) {
		return &RangeError{Key: formKey, message: fmt.Sprintf("This field cannot be more than %d characters long.", mx)}
	}

	return nil
}
