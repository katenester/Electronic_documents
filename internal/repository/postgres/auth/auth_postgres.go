package auth

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/katenester/doc/internal/models"
	"github.com/katenester/doc/internal/repository/postgres/config"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

// Создание нового пользователя (регистрация)
func (a *AuthPostgres) CreateUser(user models.User) error {
	query := fmt.Sprintf("INSERT INTO %s (login, password_hash) VALUES($1,$2)", config.UsersTable)
	_, err := a.db.Exec(query, user.Login, user.Password)
	return err
}

// Аутентификация
func (a *AuthPostgres) GetUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password_hash=$2 RETURNING id", config.UsersTable)
	row := a.db.QueryRow(query, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// Функция для получения ID пользователя по токену
func (a *AuthPostgres) GetUserId(token string) (int, error) {
	// Запрос для получения user_id по токену из таблицы sessions
	var userId int
	query := `
		SELECT user_id 
		FROM sessions 
		WHERE token = $1 AND expired_at IS NULL OR expired_at > NOW()`

	// Выполняем запрос и получаем user_id
	err := a.db.Get(&userId, query, token)
	if err != nil {
		// В случае ошибки возвращаем ошибку с пояснением
		return 0, fmt.Errorf("session not found or expired: %v", err)
	}

	// Возвращаем найденный user_id
	return userId, nil
}

// Сохранение токена для пользователя
func (a *AuthPostgres) SaveToken(userID int, token string) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, token) VALUES ($1, $2)", config.SessionsTable)
	_, err := a.db.Exec(query, userID, token)
	return err
}

// Удаление токена (завершение сессии)
func (a *AuthPostgres) DeleteToken(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = $1", config.SessionsTable)
	_, err := a.db.Exec(query, token)
	return err
}
