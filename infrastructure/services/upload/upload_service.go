package upload

import (
	"mime/multipart"
)

type UploadService interface {
	Upload(multipart.File) (string, error)
}
