package services

import (
	"crescendo-api/models"
	"crescendo-api/repositories"
)

type UserService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	service := UserService{
		repository: repository,
	}

	return service
}

func (s UserService) Create(user models.User) (int, error) {
	id, err := s.repository.Create(user)
	return id, err
}
