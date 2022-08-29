package auth

import (
	"go-boilerplate/repository/auth/entity"
	"go-boilerplate/service/auth/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(u model.User) (model.User, error) {
	user := entity.NewUser(u)
	err := r.db.Create(&user).Error
	return user.ToModel(), err
}

func (r *repository) FindOne(u model.User) (model.User, error) {
	user := entity.NewUser(u)
	err := r.db.Where(&user).First(&user).Error
	return user.ToModel(), err
}
