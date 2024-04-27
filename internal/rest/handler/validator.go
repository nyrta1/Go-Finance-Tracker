package handler

import "github.com/go-playground/validator/v10"

func validate(data any) error {
	validate := validator.New()
	return validate.Struct(data)
}
