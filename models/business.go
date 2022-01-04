package models

import (
	"gorm.io/gorm"
)

// Business for model
type Business struct {
	gorm.Model
	Name        string
	Description string
	UserID      uint
	Companies   []*Company `gorm:"many2many:businesses_companies;"`
	Articles    []*Article
	Csvs        []Csv
}

// CreateBusiness create article
func CreateBusiness(userID uint, name, description string) (*Business, string, bool) {
	var business Business
	business.UserID = userID
	business.Name = name
	business.Description = description
	result := db.Create(&business)
	if result.RowsAffected == 1 {
		return &business, "ok", true
	}
	return &business, result.Error.Error(), false
}
