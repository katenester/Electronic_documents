package service

import (
	"github.com/katenester/doc/internal/models"
	"github.com/katenester/doc/internal/repository"
	"mime/multipart"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) error
	GetUser(user models.User) (int, error)
	GetUserId(token string) (int, error)
	SaveToken(userID int, token string) error
	DeleteToken(token string) error
}

type Document interface {
	Create(fileHeader *multipart.FileHeader, doc models.Document, users []models.User) error
	GetFile(idUser int, idFile int) (*multipart.FileHeader, models.Document, error)
	GetAllFile(idUser int) ([]*multipart.FileHeader, []models.Document, error)
	DeleteFile(idUser int, idFile int) error
}

type Service struct {
	Authorization
	Document
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Document:      NewDocumentService(repos.Document),
	}
}
