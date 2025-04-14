package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
)

type S3Repo struct {
	s3Session  *minio.Client
	bucketName string
	publicUrl  string
}

func NewS3Repo(s3Session *minio.Client, bucketName string, publicUrl string) *S3Repo {
	return &S3Repo{s3Session, bucketName, publicUrl}
}

func (r *S3Repo) Upload(file multipart.File, fileName string) (string, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	_, err = r.s3Session.PutObject(
		context.Background(),
		r.bucketName,
		fileName,
		bytes.NewReader(fileBytes),
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	return r.publicUrl + r.bucketName + "/" + fileName, nil
}

func (r *S3Repo) Delete(fileName string) error {
	err := r.s3Session.RemoveObject(context.Background(), r.bucketName, fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}
