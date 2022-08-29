package auth

import (
	"go-boilerplate/service/auth/model"
	"go-boilerplate/service/auth/repository"
)

type Service interface {
	Register(u model.User) (model.User, error)
	Login(username string, password string) (model.User, error)
}

type service struct {
	repository repository.Auth
}

func NewService(repository repository.Auth) *service {
	return &service{repository}
}

func (s *service) Register(u model.User) (model.User, error) {
	newUser, err := s.repository.Create(u)
	return newUser, err
}

func (s *service) Login(username string, password string) (model.User, error) {
	u := model.User{
		Username: username,
		Password: password,
	}
	user, err := s.repository.FindOne(u)
	return user, err
}
