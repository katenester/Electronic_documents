package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/katenester/doc/internal/models"
	"net/http"
)

// Логика регистрации нового пользователя
func (h *Handler) register(c *gin.Context) {
	// Структура для запроса
	var req struct {
		Token string `json:"token" binding:"required"`
		Login string `json:"login" binding:"required"`
		PSWD  string `json:"pswd" binding:"required"`
	}

	// Получаем данные из запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrorResponse{
				Code: 400,
				Text: "Invalid parameters",
			},
		})
		return
	}

	// Проверяем токен администратора
	if req.Token != "admin_token" { // Можно извлечь токен из конфига
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": ErrorResponse{
				Code: 401,
				Text: "Not authorized",
			},
		})
		return
	}

	// Создаем нового пользователя
	user := models.User{
		Login:    req.Login,
		Password: req.PSWD,
	}

	// Сохраняем пользователя в базе
	err := h.service.Authorization.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrorResponse{
				Code: 500,
				Text: "Failed to create user",
			},
		})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"response": SuccessResponse{
			Login: req.Login,
		},
	})
}

// Логика аутентификации пользователя (вернуть токен)
func (h *Handler) signIn(c *gin.Context) {
	var req struct {
		Login string `json:"login" binding:"required"`
		PSWD  string `json:"pswd" binding:"required"`
	}

	// Получаем данные из запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrorResponse{
				Code: 400,
				Text: "Invalid parameters",
			},
		})
		return
	}
	user := models.User{
		Login:    req.Login,
		Password: req.PSWD,
	}
	// Получаем пользователя из базы данных по логину
	token, err := h.service.Authorization.GetUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": ErrorResponse{
				Code: 401,
				Text: "Invalid login or password",
			},
		})
		return
	}
	// Возвращаем токен
	c.JSON(http.StatusOK, gin.H{
		"response": token,
	})

}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Логика завершения сессии пользователя
func (h *Handler) signOut(c *gin.Context) {
	// Получаем токен из параметра URL
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrorResponse{
				Code: 400,
				Text: "Bad Request: Token is required in URL",
			},
		})
		return
	}

	// Удаляем сессию из базы данных
	err := h.service.Authorization.DeleteToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrorResponse{
				Code: 500,
				Text: "Failed to log out: " + err.Error(),
			},
		})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			token: true, // Токен, который был удален
		},
	})
}
