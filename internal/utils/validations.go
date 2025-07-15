package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
)

func ValidateParams(req interface{}) []error {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, trans)
	err := validate.Struct(req)

	if err != nil {
		var errs []error
		var validatorErrs validator.ValidationErrors
		errors.As(err, &validatorErrs)
		for _, e := range validatorErrs {
			translatedErr := fmt.Errorf(e.Translate(trans))
			errs = append(errs, translatedErr)
		}
		return errs
	}

	return nil
}
