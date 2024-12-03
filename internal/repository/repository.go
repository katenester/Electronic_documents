package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/katenester/doc/internal/models"
	"github.com/katenester/doc/internal/repository/postgres/auth"
	"github.com/katenester/doc/internal/repository/postgres/documents"
	"mime/multipart"
)

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

type Repository struct {
	Authorization
	Document
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: auth.NewAuthPostgres(db),
		Document:      documents.NewDocumentPostgres(db),
	}
}
