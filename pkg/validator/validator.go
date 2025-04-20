package validator

import (
	validatorPkg "github.com/go-playground/validator/v10"
)

var validatePkg *validatorPkg.Validate

func Get() *validatorPkg.Validate {
	return validatePkg
}

func Set() {
	validatePkg = validatorPkg.New()
}
