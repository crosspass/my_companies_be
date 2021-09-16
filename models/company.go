package models

import "gorm.io/gorm"

// Company for compnay
type Company struct {
	gorm.Model
	Name     string
	Code     string
	Profits  []Profit
	Articles []*Article `gorm:"many2many:articles_companies;"`
	Csvs     []Csv
}
