package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
)

// SearchCompany with code or name
func SearchCompany(ctx *gin.Context) {
	key := ctx.Query("key")
	var companies []models.Company
	db.Where("name like ? or code like ?", "%"+key+"%", "%"+key+"%").Find(&companies)
	ctx.JSON(http.StatusOK, gin.H{
		"code":      200,
		"companies": companies,
	})
}

// CompanyArticles for company
func CompanyArticles(ctx *gin.Context) {
	var articles = make([]models.Article, 0)
	var session models.Session
	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	var company models.Company
	db.Where("code = ?", ctx.Param("code")).Find(&company)
	db.Model(&company).Where("user_id = ?", session.UserID).Preload("Companies").Association("Articles").Find(&articles)
	ctx.JSON(http.StatusOK, gin.H{
		"code":     200,
		"articles": articles,
	})
}
