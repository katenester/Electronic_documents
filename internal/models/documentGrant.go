package models

type DocumentGrant struct {
	ID         int `json:"id"`          // Идентификатор доступа
	DocumentID int `json:"document_id"` // Идентификатор документа (ссылка на документ)
	GrantedTo  int `json:"granted_to"`  // Идентификатор пользователя, которому предоставлен доступ
}
