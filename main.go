package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/controllers"
	"github.com/my-companies-be/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Company for compnay
type Company struct {
	gorm.Model
	Name     string
	Code     string
	Profits  []Profit
	Comments []Comment
}

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

// Comment for chart
type Comment struct {
	gorm.Model
	Chart     string
	Content   string
	CompanyID uint `form:"company_id"`
	UserID    uint
}

var dsn = "host=localhost user=wu password=gorm dbname=my_companies port=5432 sslmode=disable TimeZone=Asia/Shanghai"
var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

func sliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}

//
// GET /profits?companies=1+2
func profits(c *gin.Context) {
	companies := c.Query("companies")
	companyIds := strings.Split(companies, " ")
	ids, err := sliceAtoi(companyIds)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err,
		})
	} else {
		var profits []Profit
		db.Find(&profits, ids) // find product with integer primary key
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    companyIds,
			"profit":  profits,
		})
	}
}

//
// GET /companies/sz000325
func company(c *gin.Context) {
	var company Company
	code := c.Param("code")
	fmt.Println("code", code)
	db.Where("code = ?", code).Find(&company) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"company": company,
	})
}

//
// GET /companies
func companies(c *gin.Context) {
	var companies []Company
	db.Limit(10).Find(&companies) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"companies": companies,
	})
}

//
// GET /reportSummary?code=SH600519
func reportSummary(c *gin.Context) {
	code := c.Query("code")
	var reportSummaries []models.ReportSummary
	db.Order("report_date asc").Where("company_code = ?", code).Find(&reportSummaries) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message":         "ok",
		"reportSummaries": reportSummaries,
	})
}

//
// GET /incomes?code=SH600519
func incomes(c *gin.Context) {
	code := c.Query("code")
	var incomes []models.Income
	db.Order("report_date asc").Where("company_code = ?", code).Find(&incomes) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"incomes": incomes,
	})
}

//
// GET /cashFlows?code=SH600519
func cashFlows(c *gin.Context) {
	code := c.Query("code")
	var cashFlows []models.CashFlow
	db.Order("report_date asc").Where("company_code = ?", code).Find(&cashFlows) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"cashFlows": cashFlows,
	})
}

//
// GET /balances?code=SH600519
func balances(c *gin.Context) {
	code := c.Query("code")
	var balances []models.Balance
	db.Order("report_date asc").Where("company_code = ?", code).Find(&balances)
	c.JSON(http.StatusOK, gin.H{
		"message":  "ok",
		"balances": balances,
	})
}

//
// GET /profits?companies=1+2
func saveComment(c *gin.Context) {
	var comment Comment
	err = c.BindJSON(&comment)
	log.Println("comment", comment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		result := db.Create(&comment) // pass pointer of data to Cre
		if result.RowsAffected == 1 {
			c.JSON(http.StatusOK, gin.H{
				"message":   "ok",
				"companies": comment,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": result.Error,
			})
		}
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/profits", profits)
	r.GET("/companies", companies)
	r.GET("/companies/:code", company)
	r.POST("/comments", saveComment)
	r.GET("/reportSummary", reportSummary)
	r.GET("/incomes", incomes)
	r.GET("/cashFlows", cashFlows)
	r.GET("/balances", balances)
	r.POST("/users/register", controllers.RegisterUser)
	r.POST("/users/login", controllers.Login)
	r.GET("/users/active", controllers.ActiveUser)
	r.POST("/articles", controllers.CreateArticle)
	r.GET("/articles", controllers.ListArticles)
	r.GET("/articles/:id", controllers.Article)
	r.PUT("/articles/:id", controllers.UpdateArticle)
	r.DELETE("/articles/:id", controllers.DeleteArticle)
	return r
}

func main() {
	// Read
	var company Company
	db.Preload("Profits").First(&company, 1) // find product with integer primary key
	log.Println(company.Name)
	log.Println("Profits: ", len(company.Profits))
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
