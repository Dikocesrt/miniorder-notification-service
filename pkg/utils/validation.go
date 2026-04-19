package utils

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validateInstance *validator.Validate
var once sync.Once

func initValidator() {
	validateInstance = validator.New()
}

func ValidateStruct(s any) error {
	once.Do(initValidator)
	if err := validateInstance.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var msgs []string
		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				msgs = append(msgs, fmt.Sprintf("%s is required", field))
			case "email":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid email address", field))
			case "gt":
				msgs = append(msgs, fmt.Sprintf("%s must be greater than %s", field, e.Param()))
			case "gte":
				msgs = append(msgs, fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param()))
			case "lt":
				msgs = append(msgs, fmt.Sprintf("%s must be less than %s", field, e.Param()))
			case "lte":
				msgs = append(msgs, fmt.Sprintf("%s must be less than or equal to %s", field, e.Param()))
			case "min":
				msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", field, e.Param()))
			case "max":
				msgs = append(msgs, fmt.Sprintf("%s must be at most %s characters", field, e.Param()))
			case "oneof":
				msgs = append(msgs, fmt.Sprintf("%s must be one of [%s]", field, e.Param()))
			default:
				msgs = append(msgs, fmt.Sprintf("%s failed validation on '%s'", field, e.Tag()))
			}
		}
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}
