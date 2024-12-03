package transport

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/katenester/doc/internal/models"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) uploadDocument(c *gin.Context) {
	// Логика загрузки нового документа
	// Получаем токен из заголовков или параметров
	token := c.DefaultPostForm("token", "")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token is required",
		})
		return
	}

	// Получаем информацию о файле
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File is required: %s", err),
		})
		return
	}
	defer file.Close()

	// Получаем метаданные из JSON
	var meta struct {
		Name   string   `json:"name"`
		Mime   string   `json:"mime"`
		Public bool     `json:"public"`
		Grant  []string `json:"grant"`
		Json   JSONData `json:"json"`
	}

	if err := c.ShouldBindJSON(&meta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid JSON: %s", err),
		})
		return
	}

	// Получаем пользователя по токену
	userID, err := h.service.Authorization.GetUserId(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	// Создаем документ
	doc := models.Document{
		OwnerID:  userID,
		Name:     meta.Name,
		Mime:     meta.Mime,
		Public:   meta.Public,
		JSONData: &meta.Json,
	}

	// Сохраняем файл и метаданные в базу данных
	users := []models.User{}
	for _, login := range meta.Grant {
		// Получаем пользователей по логину
		user, err := h.service.Authorization.GetUser(login)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("User %s not found", login),
			})
			return
		}
		users = append(users, user)
	}

	err = h.service.Document.Create(fileHeader, doc, users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to upload document: %s", err),
		})
		return
	}

	// Ответ с данными
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"json": doc.JSONData,
			"file": doc.Name,
		},
	})
}

func (h *Handler) getAllDocuments(c *gin.Context) {
	// Извлекаем параметры из запроса
	token := c.DefaultQuery("token", "")
	login := c.DefaultQuery("login", "")        // Опциональный параметр для фильтрации по логину
	key := c.DefaultQuery("key", "")            // Опциональный параметр для фильтрации по ключу
	value := c.DefaultQuery("value", "")        // Значение для фильтра
	limitParam := c.DefaultQuery("limit", "10") // Ограничение на количество документов

	// Парсим параметр limit
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	// Проверка токена и получение userID
	userID, err := h.service.Authorization.GetUserId(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid token",
		})
		return
	}
	files, filesInfo, err := h.service.Document.GetAllFile(userID)

	// Формируем результат
	// Реализовать

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, response)
}

func (h *Handler) getDocumentByID(c *gin.Context) {
	// Логика получения одного документа по ID
	// Извлекаем параметры из запроса
	token := c.DefaultQuery("token", "")
	docID := c.Param("id") // ID документа из URL

	// Проверка наличия токена
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token is required",
		})
		return
	}

	// Получаем ID пользователя по токену
	userID, err := h.service.Authorization.GetUserId(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid token",
		})
		return
	}

	// Получаем документ из базы данных
	docInfo, doc, err := h.service.Document.GetFile(userID, docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Document not found",
		})
		return
	}

	// Если документ файл (file = true), то отдаем его содержимое с нужным mime
	if doc.File {
		// Прочитаем файл с диска (предполагается, что путь к файлу хранится в базе данных)
		filePath := fmt.Sprintf("/path/to/documents/%s", doc.Name) // Замените на правильный путь
		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read file",
			})
			return
		}
		defer file.Close()

		// Устанавливаем правильный MIME-тип
		c.Header("Content-Type", doc.Mime)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", doc.Name))

		// Отправляем файл клиенту
		c.File(filePath)
		return
	}

	// Если это не файл, возвращаем метаданные документа в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        doc.ID,
			"name":      doc.Name,
			"mime":      doc.Mime,
			"file":      doc.File,
			"public":    doc.Public,
			"created":   doc.CreatedAt,
			"json_data": doc.JSONData,
		},
	})
}
func (h *Handler) getDocumentByIDHead(c *gin.Context) {
}
func (h *Handler) deleteDocument(c *gin.Context) {
	// Извлекаем токен из параметров запроса
	token := c.DefaultQuery("token", "")
	// Извлекаем ID документа из пути (URL) как параметр
	doc := c.Param("id") // Получаем ID документа из пути /api/docs/<id>
	docID, _ := strconv.Atoi(doc)
	// Проверка наличия токена
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token is required",
		})
		return
	}

	// Получаем ID пользователя по токену
	userID, err := h.service.Authorization.GetUserId(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Invalid token",
		})
		return
	}

	// Вызов сервиса для удаления файла (документа)
	err = h.service.Document.DeleteFile(userID, docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete document",
		})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"success": true,
		},
	})
}
