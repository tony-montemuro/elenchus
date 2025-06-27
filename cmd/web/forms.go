package main

import "github.com/tony-montemuro/elenchus/internal/validator"

type signupForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Password2 string
	validator.Validator
}

func (f signupForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["firstName"] = f.FirstName
	vals["lastName"] = f.LastName
	vals["email"] = f.Email
	vals["password"] = f.Password

	return vals
}

type loginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (f loginForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["email"] = f.Email
	vals["password"] = f.Password

	return vals
}
