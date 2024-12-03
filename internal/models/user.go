package models

import "time"

type User struct {
	ID        int       `json:"id"`         // Идентификатор пользователя
	Login     string    `json:"login"`      // Логин пользователя
	Password  string    `json:"-"`          // Пароль (хранится как хэш, не передается в JSON)
	CreatedAt time.Time `json:"created_at"` // Дата создания пользователя
}
