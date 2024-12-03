package documents

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/katenester/doc/internal/models"
	"mime/multipart"
	"os"
	"path/filepath"
)

type DocumentPostgres struct {
	db *sqlx.DB
}

func NewDocumentPostgres(db *sqlx.DB) *DocumentPostgres {
	return &DocumentPostgres{db: db}
}

// Функция для создания документа в базе данных с транзакцией
func (d *DocumentPostgres) Create(fileHeader *multipart.FileHeader, doc models.Document, users []models.User) error {
	// Генерация уникального имени для файла
	fileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)

	// Путь для сохранения файла
	savePath := "./uploads/" + fileName

	// Сохраняем файл на диск
	if err := saveFile(fileHeader, savePath); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	// Начинаем транзакцию
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			// Удаляем файл, если не удалось начать транзакцию
			deleteFile(savePath)
			tx.Rollback()
		}
	}()

	// Запись документа в таблицу `documents`
	query := `
		INSERT INTO documents (owner_id, name, mime, file, public, json_data) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`
	var docID int
	err = tx.QueryRow(query, doc.OwnerID, doc.Name, doc.Mime, true, doc.Public, doc.JSONData).Scan(&docID)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	// Сохраняем доступ для пользователей в таблице `document_grants`
	for _, user := range users {
		grantQuery := `
			INSERT INTO document_grants (document_id, granted_to) 
			VALUES ($1, $2)`
		_, err := tx.Exec(grantQuery, docID, user.ID)
		if err != nil {
			return fmt.Errorf("failed to insert document grant: %v", err)
		}
	}

	// Подтверждаем транзакцию, если все прошло успешно
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Функция для сохранения файла на диск
func saveFile(fileHeader *multipart.FileHeader, savePath string) error {
	// Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	// Создаем директорию, если ее нет
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}

	// Создаем файл на диске
	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("cannot create file: %v", err)
	}
	defer out.Close()

	// Копируем содержимое
	if _, err = file.Seek(0, 0); err != nil {
		return fmt.Errorf("cannot seek file: %v", err)
	}
	if _, err = out.ReadFrom(file); err != nil {
		return fmt.Errorf("cannot write file: %v", err)
	}

	return nil
}

// Функция для удаления файла с диска
func deleteFile(savePath string) error {
	err := os.Remove(savePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// Функция для получения одного файла
func (d *DocumentPostgres) GetFile(idUser int, idFile int) (*multipart.FileHeader, models.Document, error) {
	// Шаг 1: Проверка наличия документа в базе данных
	var doc models.Document
	query := `
		SELECT id, owner_id, name, mime, file, public 
		FROM documents 
		WHERE id = $1`
	err := d.db.Get(&doc, query, idFile)
	if err != nil {
		// Если документ не найден, возвращаем ошибку
		return nil, models.Document{}, fmt.Errorf("document not found: %v", err)
	}

	// Шаг 2: Проверка прав доступа пользователя
	if doc.Public == false {
		// Проверка наличия записи в таблице доступа (document_grants)
		var accessCheck int
		accessQuery := `
			SELECT 1 FROM document_grants 
			WHERE document_id = $1 AND granted_to = $2`
		err = d.db.Get(&accessCheck, accessQuery, idFile, idUser)
		if err != nil || accessCheck == 0 {
			// Если доступ запрещен
			return nil, models.Document{}, fmt.Errorf("user does not have access to this document")
		}
	}

	// Шаг 3: Получение файла с диска
	filePath := filepath.Join("./uploads", doc.Name) // Предположим, что файлы хранятся в каталоге "uploads"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, models.Document{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Шаг 4: Создание multipart.FileHeader для файла
	fileHeader := &multipart.FileHeader{
		Filename: doc.Name,
		Size:     int64(len(doc.Name)),      // Для примера, размер файла, если его нужно вычислить другим способом, используйте len(file)
		Header:   make(map[string][]string), // Дополнительные заголовки для файла
	}

	// Здесь можно вставить код для добавления других данных в header, если нужно, например, mime-тип.

	// Шаг 5: Возвращаем метаданные документа и сам файл
	return fileHeader, doc, nil
}

// Функция для получения всех файлов
func (d *DocumentPostgres) GetAllFile(idUser int) ([]*multipart.FileHeader, []models.Document, error) {
	// Шаг 1: Получение всех документов пользователя и доступных ему файлов
	var documents []models.Document
	query := `
		SELECT id, owner_id, name, mime, file, public 
		FROM documents
		WHERE owner_id = $1 OR id IN (SELECT document_id FROM document_grants WHERE granted_to = $1)
		ORDER BY name, created_at` // Сортировка по имени и дате создания

	err := d.db.Select(&documents, query, idUser)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving documents: %v", err)
	}

	// Шаг 2: Для каждого документа получить файл
	var fileHeaders []*multipart.FileHeader
	for _, doc := range documents {
		// Проверка, есть ли файл (file = true)
		if doc.File {
			filePath := filepath.Join("./uploads", doc.Name) // Путь к файлу на диске
			file, err := os.Open(filePath)
			if err != nil {
				// Если файл не найден, продолжим с другим документом
				continue
			}
			defer file.Close()

			// Создание multipart.FileHeader
			fileHeader := &multipart.FileHeader{
				Filename: doc.Name,
				Size:     int64(len(doc.Name)),      // Размер файла, это можно заменить на реальный размер
				Header:   make(map[string][]string), // Дополнительные заголовки для файла
			}

			// Добавляем в срез заголовков файлов
			fileHeaders = append(fileHeaders, fileHeader)
		}
	}

	// Шаг 3: Возвращаем список файлов и документов
	return fileHeaders, documents, nil
}

// Функция для удаления файла
func (d *DocumentPostgres) DeleteFile(idUser int, idFile int) error {
	// Шаг 1: Получение документа из базы данных
	var doc models.Document
	query := `
		SELECT id, owner_id, name 
		FROM documents
		WHERE id = $1`
	err := d.db.Get(&doc, query, idFile)
	if err != nil {
		return fmt.Errorf("error retrieving document: %v", err)
	}

	// Шаг 2: Проверка, является ли пользователь владельцем файла
	if doc.OwnerID != idUser {
		return fmt.Errorf("user does not have permission to delete this file")
	}

	// Шаг 3: Удаление файла с диска
	filePath := filepath.Join("./uploads", doc.Name) // Путь к файлу на диске
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error removing file: %v", err)
	}

	// Шаг 4: Удаление записи о документе из базы данных
	deleteQuery := `DELETE FROM documents WHERE id = $1`
	_, err = d.db.Exec(deleteQuery, idFile)
	if err != nil {
		return fmt.Errorf("error deleting document record: %v", err)
	}

	// Шаг 5: Удаление всех записей о доступах к документу из таблицы document_grants
	deleteGrantsQuery := `DELETE FROM document_grants WHERE document_id = $1`
	_, err = d.db.Exec(deleteGrantsQuery, idFile)
	if err != nil {
		return fmt.Errorf("error deleting document grants: %v", err)
	}

	// Успешное выполнение
	return nil
}
