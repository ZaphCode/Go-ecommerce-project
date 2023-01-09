package validation

import "github.com/go-playground/validator/v10"

type ValidationService interface {
	Validate(any) error
}

func NewValidationService() ValidationService {
	return &validationServiceImpl{validator.New()}
}

type validationServiceImpl struct {
	validator *validator.Validate
}

func (s *validationServiceImpl) Validate(data any) error {
	return s.validator.Struct(data)
}

/*
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {

	case "required":
		return "This field is required"

	case "lte":
		return "Should be less than " + fe.Param()

	case "gte":
		return "Should be greater than " + fe.Param()

	case "max":
		return "Should be less than " + fe.Param() + " characters"

	case "min":
		return "Should be greater than " + fe.Param() + " characters"

	case "email":
		return "Invalid email"
	}

	return "Unknown error"
}

func ValidateBody(body any) []*FieldError {
	var errors []*FieldError
	err := validatorInstance.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := FieldError{err.Field(), GetErrorMsg(err)}
			errors = append(errors, &element)
		}
	}
	return errors
}
*/
