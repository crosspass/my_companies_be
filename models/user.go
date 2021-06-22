package models

import (
	"errors"
	"time"

	"github.com/my-companies-be/utils"

	"gorm.io/gorm"
)

// User for model
type User struct {
	gorm.Model
	UserName      string
	FullName      string
	Email         string
	RegisterToken string
	ActiveTime    time.Time
	PasswordHash  string
	PasswordSalt  string
	IsActived     bool
}

// Active user
func (u *User) Active(tx *gorm.DB) {
	tx.Model(u).Update("is_actived", true)
}

// BeforeCreate for hook validate
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserName == "" || u.Email == "" {
		err = errors.New("username or email can not be blank")
		return
	}
	var count int64
	tx.Table("users").Where("email = ? or user_name = ?", u.Email, u.UserName).Count(&count)
	if count != 0 {
		err = errors.New("email has register")
		return
	}
	token, err := utils.GenerateRandomString(20)
	for tx.Table("users").Where("register_token = ?", token).Count(&count); count != 0; {
		token, err = utils.GenerateRandomString(20)
	}
	u.RegisterToken = token
	return
}
