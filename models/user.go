package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
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
	Session       Session
	Companies     []Company `gorm:"many2many:users_companies;"`
}

// Active user
func (u *User) Active(tx *gorm.DB) {
	tx.Model(u).Update("is_actived", true)
}

// SetPassword set user password
func (u *User) SetPassword(password string) {
	sum := sha256.Sum256([]byte(password + u.PasswordSalt))
	hash := fmt.Sprintf("%x", sum)
	db.Model(&u).Update("password_hash", hash)
}

// ValidatePassword validate user login by password
func (u *User) ValidatePassword(password string) bool {
	sum := sha256.Sum256([]byte(password + u.PasswordSalt))
	hash := fmt.Sprintf("%x", sum)
	return u.PasswordHash == hash
}

// GenerateToken user login successful generate token for api auth
func (u *User) GenerateToken() string {
	token, err := utils.GenerateRandomString(20)
	if err != nil {
		log.Panic(err)
	}
	var count int64
	for db.Table("sessions").Where("key = ?", token).Count(&count); count != 0; {
		token, err = utils.GenerateRandomString(20)
		if err != nil {
			log.Panic(err)
		}
	}
	if u.Session.Key == "" {
		db.Model(u).Association("Session").Append(&Session{Key: token, LoginTime: time.Now()})
	} else {
		u.Session.Key = token
		u.Session.LoginTime = time.Now()
		db.Save(&u.Session)
	}
	return token
}

// BeforeCreate for hook validate
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Email == "" {
		err = errors.New("username or email can not be blank")
		return
	}
	var count int64
	tx.Table("users").Where("email = ?", u.Email).Count(&count)
	if count != 0 {
		err = errors.New("email has register")
		return
	}
	token, err := utils.GenerateRandomString(20)
	for tx.Table("users").Where("register_token = ?", token).Count(&count); count != 0; {
		token, err = utils.GenerateRandomString(20)
	}
	u.RegisterToken = token
	u.PasswordSalt = token
	return
}
