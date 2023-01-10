package upload

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func NewFirebaseUploadService() UploadService {
	return &firebaseUploadServiceImpl{}
}

type firebaseUploadServiceImpl struct{}

func (s *firebaseUploadServiceImpl) Upload(file multipart.File) (string, error) {
	return "", nil
}

func (s *firebaseUploadServiceImpl) ValidateFileType(file multipart.File, extencions []string) error {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(bytes) // []string{"image/png", "image/jpeg", "image/jpg"}

	if ok := s.stringInSlice(mimeType, extencions); !ok {
		return fmt.Errorf("file type: %s is not supported", mimeType)
	}

	return nil
}

func (s *firebaseUploadServiceImpl) stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
