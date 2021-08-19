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
}

// CreateArticle create article
func CreateArticle(userID uint, HTMLContent, rawContent string) (string, bool) {
	var article Article
	article.UserID = userID
	article.Content = HTMLContent
	article.RawContent = rawContent
	result := db.Create(&article) // pass pointer of data to Cre
	if result.RowsAffected == 1 {
		return "ok", true
	}
	return result.Error.Error(), false
}
