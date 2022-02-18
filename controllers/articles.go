package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
)

// ArticleReq for create article form
type ArticleReq struct {
	ID          uint
	HTMLContent string `json:"htmlContent"`
	RawContent  string `json:"rawcontent"`
	CompanyIds  []int  `json:"company_ids"`
}

// CreateArticle create article
func CreateArticle(c *gin.Context) {
	var session models.Session
	var companies []models.Company
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Preload("User").Where("key = ?", token).Find(&session)
	var articleReq ArticleReq
	err := c.BindJSON(&articleReq)
	log.Println("artilceReq", articleReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		article, msg, ok := models.CreateArticle(session.UserID, articleReq.HTMLContent, articleReq.RawContent)
		if ok {
			db.Where("id IN ?", articleReq.CompanyIds).Find(&companies)
			db.Model(&article).Association("Companies").Append(companies)
			db.Model(&session.User).Association("Companies").Append(companies)
			c.JSON(http.StatusOK, gin.H{
				"article": article,
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
	var companies []models.Company
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Preload("User").Where("key = ?", token).Find(&session)
	var articleReq ArticleReq
	var article models.Article
	err := c.BindJSON(&articleReq)
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
			article.Content = articleReq.HTMLContent
			article.RawContent = articleReq.RawContent
			db.Save(&article)
			fmt.Println("company_ids: ", articleReq.CompanyIds)
			db.Where("id IN ?", articleReq.CompanyIds).Find(&companies)
			db.Model(&article).Association("Companies").Replace(companies)
			db.Model(&session.User).Association("Companies").Append(companies)
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		}
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
	// page := c.DefaultQuery("page", "1")
	// offset, error := strconv.Atoi(page)
	// if error != nil {
	// 	log.Fatal("page format error", error)
	// }
	db.Preload("Companies").Where("user_id = ? AND date_part('year',created_at) = ?", session.UserID, year).Find(&articles) //.Offset(offset * 20).Limit(20)
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
	db.Preload("Companies").Where("user_id = ? AND ID = ?", session.UserID, id).Find(&article)
	c.JSON(http.StatusOK, gin.H{
		"article": article,
		"message": "ok",
	})
}

type statsStruct struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// StatsArticle get specified article by ID
func StatsArticle(c *gin.Context) {
	var session models.Session
	token := c.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	year := c.Query("year")
	start := year + "-01-01"
	end := year + "-12-31"
	sqlTmp := `SELECT date(created_at) as date, count(created_at) as total_count 
	FROM "articles"
	WHERE (deleted_at is null AND "articles"."user_id" = %d AND "articles"."created_at" BETWEEN '%s' AND '%s') 
	GROUP BY date(created_at)
	ORDER BY date DESC`
	sql := fmt.Sprintf(sqlTmp, session.UserID, start, end)
	log.Printf("sql %s\n", sql)
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		log.Printf("sql %s, error: %s", sql, err)
	}
	defer rows.Close()
	var date string
	var count int
	var stats []statsStruct
	for rows.Next() {
		rows.Scan(&date, &count)
		log.Printf("date: %s, count: %d\n", date, count)
		stats = append(stats, statsStruct{date, count})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"stats":   stats,
	})
}

// DeleteArticle mark article as deleted
func DeleteArticle(c *gin.Context) {
	var session models.Session
	var article models.Article
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Where("key = ?", token).Find(&session)
	id := c.Param("id")
	db.Where("user_id = ? AND ID = ?", session.UserID, id).Find(&article)
	db.Delete(&article)
	c.JSON(http.StatusOK, gin.H{
		"article": article,
		"message": "ok",
	})
}
