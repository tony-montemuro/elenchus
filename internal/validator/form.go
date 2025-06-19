package validator

type Form interface {
	ValidatorInterface
	GetStringVals() map[string]string
}

var (
	SignUpForm = "signupForm"
)
