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

// ReportSummary for company
type ReportSummary struct {
	gorm.Model
	Category                            string
	CompanyCode                         string
	ReportName                          string
	ReportDate                          uint64
	AvgRoe                              float32
	AvgRoeIncrease                      float32
	NpPerShare                          float32
	NpPerShareIncrease                  float32
	OperateCashFlowPs                   float32
	OperateCashFlowPsIncrease           float32
	BasicEps                            float32
	BasicEpsIncrease                    float32
	CapitalReserve                      float32
	CapitalReserveIncrease              float32
	UndistriProfitPs                    float32
	UndistriProfitPsIncrease            float32
	NetInterestOfTotalAssets            float32
	NetInterestOfTotalAssetsIncrease    float32
	NetSellingRate                      float32
	NetSellingRateIncrease              float32
	GrossSellingRate                    float32
	GrossSellingRateIncrease            float32
	TotalRevenue                        float32
	TotalRevenueIncrease                float32
	OperatingIncomeYoy                  float32
	OperatingIncomeYoyIncrease          float32
	NetProfitAtsopc                     float32
	NetProfitAtsopcIncrease             float32
	NetProfitAtsopcYoy                  float32
	NetProfitAtsopcYoyIncrease          float32
	NetProfitAfterNrgalAtsolc           float32
	NetProfitAfterNrgalAtsolcIncrease   float32
	NpAtsopcNrgalYoy                    float32
	NpAtsopcNrgalYoyIncrease            float32
	OreDlt                              float32
	OreDltIncrease                      float32
	Rop                                 float32
	RopIncrease                         float32
	AssetLiabRatio                      float32
	AssetLiabRatioIncrease              float32
	CurrentRatio                        float32
	CurrentRatioIncrease                float32
	QuickRatio                          float32
	QuickRatioIncrease                  float32
	EquityMultiplier                    float32
	EquityMultiplierIncrease            float32
	EquityRatio                         float32
	EquityRatioIncrease                 float32
	HolderEquity                        float32
	HolderEquityIncrease                float32
	NcfFromOaToTotalLiab                float32
	NcfFromOaToTotalLiabIncrease        float32
	InventoryTurnoverDays               float32
	InventoryTurnoverDaysIncrease       float32
	ReceivableTurnoverDays              float32
	ReceivableTurnoverDaysIncrease      float32
	AccountsPayableTurnoverDays         float32
	AccountsPayableTurnoverDaysIncrease float32
	CashCycle                           float32
	CashCycleIncrease                   float32
	OperatingCycle                      float32
	OperatingCycleIncrease              float32
	TotalCapitalTurnover                float32
	TotalCapitalTurnoverIncrease        float32
	InventoryTurnover                   float32
	InventoryTurnoverIncrease           float32
	AccountReceivableTurnover           float32
	AccountReceivableTurnoverIncrease   float32
	AccountsPayableTurnover             float32
	AccountsPayableTurnoverIncrease     float32
	CurrentAssetTurnoverRate            float32
	CurrentAssetTurnoverRateIncrease    float32
	FixedAssetTurnoverRatio             float32
	FixedAssetTurnoverRatioIncrease     float32
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
// GET /profits?companies=1+2
func company(c *gin.Context) {
	var company Company
	db.Preload("Profits").Preload("Comments").First(&company, 1) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"company": company,
	})
}

//
// GET /companies
func companies(c *gin.Context) {
	var companies []Company
	db.Find(&companies) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"companies": companies,
	})
}

//
// GET /reportSummary?company_id=1
func reportSummary(c *gin.Context) {
	code := c.Query("code")
	var reportSummaries []ReportSummary
	db.Where("company_code = ?", code).Find(&reportSummaries) // find product with integer primary key
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"company": reportSummaries,
	})
}

//
// GET /profits?companies=1+2
func saveComment(c *gin.Context) {
	var comment Comment
	c.BindJSON(&comment)
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
	r.GET("/companies/:id", company)
	r.POST("/comments", saveComment)
	r.GET("/reportSummary", reportSummary)
	return r
}

func main() {
	// Read
	var company Company
	db.Preload("Profits").First(&company, 1) // find product with integer primary key
	log.Println(company.Name)
	log.Println("Profits: ", len(company.Profits))
	// var profit Profit
	// db.First(&profit, 1) // find product with integer primary key
	// log.Println(profit.Ying_Shou)
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
