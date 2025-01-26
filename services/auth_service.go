package services

import (
	"gin-fleamarket/models"
	"gin-fleamarket/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(email string, password string) error
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Signup(email string, password string) error {
	// パスワードのハッシュ化
	// ハッシュ化のサンプル https://pkg.go.dev/golang.org/x/crypto/bcrypt#GenerateFromPassword
	// ハッシュ化のサンプル https://pkg.go.dev/golang.org/x/crypto/bcrypt#CompareHashAndPassword
	// ハッシュ化のサンプル https://go.dev/play/p/8t7j7r5vL2w
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repository.CreateUser(user)
}
