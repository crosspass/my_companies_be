package models

import "gorm.io/gorm"

// Profit is for company's profit
type Profit struct {
	gorm.Model
	Year           string
	YingShou       int64
	YingYeChengBen int64
	FeiYingShou    int64
	LiRun          int64
	YingLiRun      int64
	CompanyID      uint
}
