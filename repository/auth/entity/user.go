package entity

import (
	"go-boilerplate/service/auth/model"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(m model.User) User {
	return User{
		ID:       m.ID,
		Username: m.Username,
		Password: m.Password,
		Email:    m.Email,
		Name:     m.Name,
	}
}

func (u User) ToModel() model.User {
	return model.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Name:     u.Name,
	}
}
