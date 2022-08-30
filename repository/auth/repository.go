package auth

import (
	"errors"
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

func (r *repository) Register(u model.User) (model.User, error) {
	var us entity.User
	user := entity.NewUser(u)
	r.db.Where(entity.User{
		Username: u.Username,
	}).
		Or(entity.User{
			Email: u.Email,
		}).
		First(&us)

	if us.ID > 0 {
		return us.ToModel(), errors.New("Data telah dibuat")
	} else {
		err := r.db.Create(&user).Error
		return user.ToModel(), err
	}
	// err := r.db.Create(&user).Error
}
