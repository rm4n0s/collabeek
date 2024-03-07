package common

import (
	"encoding/json"

	"github.com/go-playground/validator"
)

type ValidatorErrorJson struct {
	Field   string
	Type    string
	Message string
}

type ListValidatorErrorJson []ValidatorErrorJson

func (e *ListValidatorErrorJson) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		var vers ListValidatorErrorJson
		for _, e := range errs {
			vers = append(vers, ValidatorErrorJson{
				Field:   e.Field(),
				Type:    e.Type().Name(),
				Message: e.Translate(nil),
			})
		}
		return &vers
	}
	return nil
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}
