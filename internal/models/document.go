package models

import "time"

type Document struct {
	ID        int       `json:"id"`                  // Идентификатор документа
	OwnerID   int       `json:"owner_id"`            // Идентификатор владельца (ссылка на пользователя)
	Name      string    `json:"name"`                // Имя документа
	Mime      string    `json:"mime"`                // MIME-тип документа
	File      bool      `json:"file"`                // Флаг наличия файла
	Public    bool      `json:"public"`              // Флаг публичности документа
	JSONData  *JSONData `json:"json_data,omitempty"` // Метаданные документа в формате JSON (может отсутствовать)
	CreatedAt time.Time `json:"created_at"`          // Дата создания документа
	UpdatedAt time.Time `json:"updated_at"`          // Дата обновления документа
}

// Структура для метаданных документа, если они есть
type JSONData map[string]interface{}
