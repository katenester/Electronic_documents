package models

import "time"

type Session struct {
	ID        int       `json:"id"`         // Идентификатор сессии
	UserID    int       `json:"user_id"`    // Идентификатор пользователя (ссылка на пользователя)
	Token     string    `json:"token"`      // Токен сессии
	CreatedAt time.Time `json:"created_at"` // Дата создания сессии
	ExpiredAt time.Time `json:"expired_at"` // Дата истечения срока действия сессии
}
