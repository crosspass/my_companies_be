package models

import (
	"time"

	"gorm.io/gorm"
)

// Session for model
type Session struct {
	gorm.Model
	Key       string
	LoginTime time.Time
	UserID    uint
	User      *User
}
