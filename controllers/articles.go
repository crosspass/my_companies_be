package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
)

// ArticleReq for create article form
type ArticleReq struct {
	ID          uint
	HTMLContent string `json:"htmlContent"`
	RawContent  string `json:"rawcontent"`
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
		msg, ok := models.CreateArticle(session.UserID, articleReq.HTMLContent, articleReq.RawContent)
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

// UpdateArticle create article
func UpdateArticle(c *gin.Context) {
	var session models.Session
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Where("key = ?", token).Find(&session)
	var articleReq ArticleReq
	var article models.Article
	err = c.BindJSON(&articleReq)
	log.Println("artilceReq", articleReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		ret := db.First(&article, articleReq.ID)
		if ret.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": ret.Error.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		}
		article.Content = articleReq.HTMLContent
		article.RawContent = articleReq.RawContent
		db.Save(&article)
	}
}

// PageSize for page records size
var PageSize = 20

// ListArticles get
func ListArticles(c *gin.Context) {
	var session models.Session
	var articles []models.Article
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Where("key = ?", token).Find(&session)
	year := c.DefaultQuery("year", time.Now().Format("2006"))
	page := c.DefaultQuery("page", "1")
	offset, error := strconv.Atoi(page)
	if error != nil {
		log.Fatal("page format error", error)
	}
	db.Where("user_id = ? AND date_part('year',created_at) = ?", session.UserID, year).Find(&articles).Offset(offset * 20).Limit(20)
	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"message":  "ok",
	})
}

// Article get specified article by ID
func Article(c *gin.Context) {
	var session models.Session
	var article models.Article
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Where("key = ?", token).Find(&session)
	id := c.Param("id")
	db.Where("user_id = ? AND ID = ?", session.UserID, id).Find(&article)
	c.JSON(http.StatusOK, gin.H{
		"article": article,
		"message": "ok",
	})
}
