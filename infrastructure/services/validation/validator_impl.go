package validation

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Custom types
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (ValidationErrors) Error() string {
	return "someting wentwrong"
}

//* Constructor
func NewValidationService() ValidationService {
	return &validationServiceImpl{validator: validator.New()}
}

//* Implementation
type validationServiceImpl struct {
	validator *validator.Validate
}

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

func (s *validationServiceImpl) ValidateFileType(file multipart.File, extencions []string) error {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(bytes) // []string{"image/png", "image/jpeg", "image/jpg"}

	if ok := s.stringInSlice(mimeType, extencions); !ok {
		return fmt.Errorf("file type: %s is not supported", mimeType)
	}

	return nil
}

func (s *validationServiceImpl) stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
	}
	return "Unknown error"
}
