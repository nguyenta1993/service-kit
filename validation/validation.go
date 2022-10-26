package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type CustomValidation struct {
	Tag           string
	ValidatorFunc validator.Func
}

func UseValidation(customValidations ...CustomValidation) {
	validate = validator.New()

	for _, customValidation := range customValidations {
		err := validate.RegisterValidation(customValidation.Tag, customValidation.ValidatorFunc)
		if err != nil {
			return
		}

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation(customValidation.Tag, customValidation.ValidatorFunc)
			if err != nil {
				return
			}
		}
	}
}
