package validator

import (
	"fmt"
	"strings"
)

type RangeRule struct {
	MinLength int
	MaxLength int
}

type FormRangeRules map[string]RangeRule

var RangeRules = map[string]FormRangeRules{
	SignUpForm: {
		"firstName": {MinLength: 1, MaxLength: 50},
		"lastName":  {MinLength: 1, MaxLength: 50},
		"email":     {MinLength: 1, MaxLength: 255},
		"password":  {MinLength: 8, MaxLength: 72},
	},
}

func InputsInRange(form Form, name string) error {
	rules := RangeRules[name]
	missed := []string{}

	for key, val := range form.GetStringVals() {
		rule, exists := rules[key]
		if exists {
			checkFields(form, rule, key, val)
		}

		if !exists {
			missed = append(missed, key)
		}
	}

	if len(missed) > 0 {
		return fmt.Errorf("The following range checks failed: %s", strings.Join(missed, ", "))
	}

	return nil
}

func checkFields(form Form, rule RangeRule, formKey, formValue string) {
	mn, mx := rule.MinLength, rule.MaxLength
	if mn > 0 {
		if mn == 1 {
			form.CheckField(NotBlank(formValue), formKey, "This field cannot be blank.")
		} else {
			form.CheckField(MinChars(formValue, mn), formKey, fmt.Sprintf("This field must be at least %d characters long.", mn))
		}
	}
	if mx > 0 {
		form.CheckField(MaxChars(formValue, mx), formKey, fmt.Sprintf("This field cannot be more than %d characters long.", mx))
	}
}
