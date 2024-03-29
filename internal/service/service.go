package service

import (
	"github.com/eydeveloper/highload-social/internal/entity"
	"github.com/eydeveloper/highload-social/internal/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (string, error)
	GenerateToken(id string, password string) (string, error)
}

type User interface {
	GetById(id string) (repository.UserProfile, error)
}

type Service struct {
	Authorization
	User
}

func NewService(repositories *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repositories.Authorization),
		User: NewUserService(repositories.User),
	}
}
