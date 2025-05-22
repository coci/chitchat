package repository

import (
	"errors"
	"github.com/coci/chitchat/pkg/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {

	return &UserRepository{}
}

func (r UserRepository) StoreUser(username, password string) (*model.User, error) {
	user := model.User{
		Username: username,
		Password: password,
	}
	result := r.db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}

func (r UserRepository) GetUser(username, password string) (*model.User, error) {
	user := model.User{
		Username: username,
	}

	result := r.db.First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if user.Password != password {
		return nil, errors.New("wrong password")
	}
	return &user, nil
}
