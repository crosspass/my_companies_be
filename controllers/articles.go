package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
)

// ArticleReq for create article form
type ArticleReq struct {
	Content string `json:"content"`
}

// CreateArticle create article
func CreateArticle(c *gin.Context) {
	var session models.Session
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Where("key = ?", token).Find(&session)
	var articleReq ArticleReq
	err = c.BindJSON(&articleReq)
	log.Println("artilceReq", articleReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		msg, ok := models.CreateArticle(session.UserID, articleReq.Content)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": msg,
			})
		}
	}
}
