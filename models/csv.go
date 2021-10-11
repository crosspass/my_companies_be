package models

import "gorm.io/gorm"

// Csv is for company's specified csv file
type Csv struct {
	gorm.Model
	UserID    uint
	CompanyID uint
	Title     string `json:"title"`
	ChartType string `json:"chartType"`
	Data      string `json:"data"`
}
