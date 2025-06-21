package validator

type Form interface {
	GetStringVals() map[string]string
}

var (
	SignUpForm = "signupForm"
	LoginForm  = "loginForm"
)
