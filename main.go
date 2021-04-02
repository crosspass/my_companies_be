package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Company for compnay
type Company struct {
	gorm.Model
	Name string
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
	CompanyID      int64
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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/profits", profits)
	return r
}

func main() {
	// Read
	var company Company
	db.First(&company, 1) // find product with integer primary key
	log.Println(company.Name)
	// var profit Profit
	// db.First(&profit, 1) // find product with integer primary key
	// log.Println(profit.Ying_Shou)
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
