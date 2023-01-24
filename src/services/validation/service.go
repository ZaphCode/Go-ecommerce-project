package validation

import (
	"mime/multipart"
)

//* Service

type ValidationService interface {
	Validate(any) error
	ValidateFileType(f multipart.File, exts ...string) error
}
