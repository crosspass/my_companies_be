package models

import "gorm.io/gorm"

// Csv is for company's specified csv file
type Csv struct {
	gorm.Model
	Path       string
	Name       string
	OriginName string
	Size       int64
	UserID     uint
	CompanyID  uint
}
