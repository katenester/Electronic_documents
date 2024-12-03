package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/katenester/doc/internal/models"
	"github.com/katenester/doc/internal/repository"
	"regexp"
	"time"
)

const (
	salt       = "sfsgGhJjJJHgFRdehYgu"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH" // Ключ подписи
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) CreateUser(user models.User) error {
	if validateLogin(user.Login) && validatePassword(user.Password) {
		user.Password = generatePasswordHash(user.Password)
		return s.repo.CreateUser(user)
	} else {
		return errors.New("Invalid username or password")
	}
}
func validatePassword(password string) bool {
	// Минимум 8 символов
	if len(password) < 8 {
		return false
	}
	// Минимум 1 заглавная и 1 строчная буква, минимум 1 цифра и 1 спецсимвол
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\d]).{8,}$`
	match, _ := regexp.MatchString(regex, password)
	return match
}

// Функция для валидации логина
func validateLogin(login string) bool {
	// Минимум 8 символов, латиница и цифры
	regex := `^[a-zA-Z0-9]{8,}$`
	match, _ := regexp.MatchString(regex, login)
	return match
}
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
func (s *AuthService) GetUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.GetUser(user)
}
func (s *AuthService) GetUserId(token string) (int, error) {
	return s.repo.GetUserId(token)
}
func (s *AuthService) SaveToken(userID int, token string) error {
	return s.repo.SaveToken(userID, token)
}
func (s *AuthService) DeleteToken(token string) error {
	return s.repo.DeleteToken(token)
}

//func (s *AuthService) GenerateToken(username, password string) (string, error) {
//	user, err := s.repo.GetUser(username, generatePasswordHash(password))
//	if err != nil {
//		return "", err
//	}
//	// Если пользователь существует => генерируем токен
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
//		jwt.StandardClaims{
//			// Дедлайн валидности токена
//			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
//			// Время создания токена
//			IssuedAt: time.Now().Unix(),
//		},
//		user.Id,
//	})
//	return token.SignedString([]byte(signingKey))
//}
//
//func (s *AuthService) ParseToken(accessToken string) (int, error) {
//	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("invalid signature method")
//		}
//		return []byte(signingKey), nil
//	})
//	if err != nil {
//		return 0, err
//	}
//	claims, ok := token.Claims.(*tokenClaims)
//	if !ok {
//		return 0, errors.New("token claims is not of type *tokenClaims")
//	}
//	return claims.UserId, nil
//}
