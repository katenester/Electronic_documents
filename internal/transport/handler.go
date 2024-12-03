package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/katenester/doc/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	// Группа для работы с аутентификацией и регистрацией
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register) // Регистрация нового пользователя
		auth.POST("/auth", h.signIn)       // Аутентификация пользователя
		auth.DELETE("/:token", h.signOut)  // Завершение авторизованной сессии
	}

	// Группа для работы с документами (защищенные маршруты)
	api := router.Group("/api")
	{
		// Работа с документами
		docs := api.Group("/docs")
		{
			docs.POST("/", h.uploadDocument)         // Загрузка нового документа
			docs.GET("/", h.getAllDocuments)         // Получение списка документов
			docs.GET("/:id", h.getDocumentByID)      // Получение одного документа
			docs.HEAD("/:id", h.getDocumentByIDHead) // HEAD запрос для документа
			docs.DELETE("/:id", h.deleteDocument)    // Удаление документа
		}
	}
	return router
}
