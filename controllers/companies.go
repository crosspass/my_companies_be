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
