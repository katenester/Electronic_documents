package service

import (
	"github.com/katenester/doc/internal/models"
	"github.com/katenester/doc/internal/repository"
	"mime/multipart"
)

type DocumentService struct {
	repo repository.Document
}

func NewDocumentService(repo repository.Document) *DocumentService {
	return &DocumentService{repo: repo}
}

func (d DocumentService) Create(fileHeader *multipart.FileHeader, doc models.Document, users []models.User) error {
	return d.repo.Create(fileHeader, doc, users)
}
func (d DocumentService) GetFile(idUser int, idFile int) (*multipart.FileHeader, models.Document, error) {
	return d.repo.GetFile(idUser, idFile)
}
func (d DocumentService) GetAllFile(idUser int) ([]*multipart.FileHeader, []models.Document, error) {
	return d.repo.GetAllFile(idUser)
}
func (d DocumentService) DeleteFile(idUser int, idFile int) error {
	return d.repo.DeleteFile(idUser, idFile)
}
