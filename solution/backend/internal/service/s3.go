package service

import (
	"mime/multipart"
	"neJok/solution/internal/repository"
)

type S3Service struct {
	S3Repo *repository.S3Repo
}

func NewS3Service(S3Repo *repository.S3Repo) *S3Service {
	return &S3Service{S3Repo}
}

func (s *S3Service) UploadToS3(file multipart.File, fileName string) (string, error) {
	return s.S3Repo.Upload(file, fileName)
}

func (s *S3Service) Delete(fileName string) error {
	return s.S3Repo.Delete(fileName)
}
