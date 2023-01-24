package validation

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/go-playground/validator/v10"
)

//* Implementation

type validationServiceImpl struct {
	validator *validator.Validate
}

//* Constructor

func NewValidationService() ValidationService {
	return &validationServiceImpl{validator: validator.New()}
}

//* Methods

func (s *validationServiceImpl) Validate(data any) error {
	err := s.validator.Struct(data)
	if err != nil {
		var errors ValidationErrors
		for _, err := range err.(validator.ValidationErrors) {
			element := FieldError{err.Field(), s.getErrorMsg(err)}
			errors = append(errors, element)
		}
		return errors
	}
	return nil
}

func (s *validationServiceImpl) ValidateFileType(f multipart.File, exts ...string) error {
	bt, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(bt) // []string{"image/png", "image/jpeg", "image/jpg"}

	if ok := utils.ItemInSlice(mimeType, exts); !ok {
		return fmt.Errorf("file type: %s is not supported", mimeType)
	}

	return nil
}

func (s *validationServiceImpl) getErrorMsg(fe validator.FieldError) string {
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
	case "oneof":
		return fmt.Sprintf("Yo can only choise between: [%s]", fe.Param())
	case "url":
		return "Invalid url"
	default:
		return fmt.Sprintf("Unknown error (%s)", fe.Tag())
	}
}

//* Custom types

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (ValidationErrors) Error() string {
	return "something wentwrong"
}
