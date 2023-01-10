package validation

import (
	"mime/multipart"
)

//* Service
type ValidationService interface {
	Validate(any) error
	ValidateFileType(file multipart.File, extencions []string) error
}
