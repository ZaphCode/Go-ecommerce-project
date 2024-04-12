package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"firebase.google.com/go/v4/storage"
	"github.com/google/uuid"
)

//* Implementation

type firebaseUploadServiceImpl struct {
	client     *storage.Client
	bucketName string
	folder     string
}

//* Constructor

func NewFirebaseUploadService(
	client *storage.Client, bucketName, folder string,
) UploadService {
	return &firebaseUploadServiceImpl{
		client:     client,
		bucketName: bucketName,
		folder:     folder,
	}
}

//* Methods

func (s *firebaseUploadServiceImpl) Upload(file multipart.File) (string, error) {
	bucket, err := s.client.Bucket(s.bucketName)

	if err != nil {
		return "", fmt.Errorf("client.Bucket() error : %w", err)
	}

	id, err := uuid.NewUUID()

	if err != nil {
		return "", fmt.Errorf("uuid.NewUUID() error : %w", err)
	}

	wc := bucket.Object(s.folder + "/" + id.String()).NewWriter(context.Background())

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media",
		s.bucketName,
		s.folder+"%2F"+id.String(),
	), nil
}
