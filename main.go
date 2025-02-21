package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/connect"
	"github.com/my-companies-be/controllers"
	controllersV2 "github.com/my-companies-be/controllers/v2"
	"github.com/my-companies-be/models"
	"github.com/spf13/viper"
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

// var dsn = "host=localhost user=wu password=gorm dbname=my_companies port=5432 sslmode=disable TimeZone=Asia/Shanghai"
// var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// var dsn = "host=localhost user=wu password=gorm dbname=my_companies port=5432 sslmode=disable TimeZone=Asia/Shanghai"
var db = connect.Db

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

// GET /profits?companies=1+2
func saveComment(c *gin.Context) {
	var comment Comment
	err := c.BindJSON(&comment)
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
	upload := viper.GetString("uploads")
	r.Static("/uploads", upload)
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/profits", profits)
	r.GET("/companies", controllers.Companies)
	r.GET("/companies/:code", company)
	r.GET("/companies/:code/articles", controllers.CompanyArticles)
	r.GET("/company/search", controllers.SearchCompany)
	r.POST("/comments", saveComment)
	r.GET("/reportSummary", reportSummary)
	r.GET("/incomes", incomes)
	r.GET("/cashFlows", cashFlows)
	r.GET("/balances", balances)
	r.POST("/users/register", controllers.RegisterUser)
	r.POST("/users/login", controllers.Login)
	r.GET("/users/active", controllers.ActiveUser)
	r.GET("/user/info", controllers.Info)
	r.POST("/users/starCompany", controllers.StarCompany)
	r.PUT("/unstarCompany", controllers.UnStarCompany)
	r.POST("/articles", controllers.CreateArticle)
	r.GET("/articles", controllers.ListArticles)
	r.GET("/articles/:id", controllers.Article)
	r.PUT("/articles/:id", controllers.UpdateArticle)
	r.DELETE("/articles/:id", controllers.DeleteArticle)
	r.GET("/article/stats", controllers.StatsArticle)
	r.GET("/csvs", controllers.IndexCsv)
	r.POST("/csvs/upload", controllers.UploadCSV)
	r.POST("/csvs/upload/:id", controllers.UpdateCSVFile)
	r.POST("/csvs", controllers.CreateCsv)
	r.PUT("/csvs/:id", controllers.UpdateCsv)
	r.POST("/businesses", controllers.CreateBusiness)
	r.GET("/businesses", controllers.ListBusiness)
	r.GET("/businesses/:id", controllers.Business)
	r.GET("/businesses/:id/stats", controllers.BusinessStats)
	r.PUT("/businesses/:id", controllers.UpdateBusiness)
	r.GET("/tops/roa", controllers.Roa)
	r.GET("/tops/roaIncrease", controllers.RoaIncrease)
	r.GET("/tops/roe", controllers.Roe)
	r.GET("/tops/roeIncrease", controllers.RoeIncrease)

	// v2 api
	r.POST("/v2/users/register", controllersV2.RegisterUser)
	r.POST("/v2/users/login", controllersV2.Login)
	r.GET("/v2/users/active", controllersV2.ActiveUser)
	r.GET("/v2/user/info", controllersV2.Info)
	r.POST("/v2/users/starCompany", controllersV2.StarCompany)
	r.PUT("/v2/unstarCompany", controllersV2.UnStarCompany)
	return r
}

func main() {
	// Read
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
