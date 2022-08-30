package repository

import "go-boilerplate/service/auth/model"

type Auth interface {
	Create(u model.User) (model.User, error)
	FindOne(u model.User) (model.User, error)
	Register(u model.User) (model.User, error)
}
