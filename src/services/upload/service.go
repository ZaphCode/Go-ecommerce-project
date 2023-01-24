package upload

import (
	"mime/multipart"
)

//* Service

type UploadService interface {
	Upload(multipart.File) (string, error)
}
