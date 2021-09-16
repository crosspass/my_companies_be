package models

import (
	"gorm.io/gorm"
)

// Article for model
type Article struct {
	gorm.Model
	Content    string
	RawContent string
	UserID     uint
	Companies  []*Company `gorm:"many2many:articles_companies;"`
}

// CreateArticle create article
func CreateArticle(userID uint, HTMLContent, rawContent string) (*Article, string, bool) {
	var article Article
	article.UserID = userID
	article.Content = HTMLContent
	article.RawContent = rawContent
	result := db.Create(&article) // pass pointer of data to Cre
	if result.RowsAffected == 1 {
		return &article, "ok", true
	}
	return &article, result.Error.Error(), false
}
