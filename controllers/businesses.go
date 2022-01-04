package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
)

// BusinessReq for create business form
type BusinessReq struct {
	ID          uint
	Name        string `json:"name"`
	Description string `json:"description"`
	CompanyIds  []int  `json:"company_ids"`
}

// CreateBusiness create business
func CreateBusiness(c *gin.Context) {
	user := currentUser(c)
	var companies []models.Company
	var businessReq BusinessReq
	err := c.BindJSON(&businessReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		business, msg, ok := models.CreateBusiness(user.ID, businessReq.Name, businessReq.Description)
		if ok {
			db.Where("id IN ?", businessReq.CompanyIds).Find(&companies)
			db.Model(&business).Association("Companies").Append(companies)
			db.Model(&user).Association("Companies").Append(companies)
			c.JSON(http.StatusOK, gin.H{
				"business": business,
				"message":  "ok",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": msg,
			})
		}
	}
}

// UpdateBusiness create business
func UpdateBusiness(c *gin.Context) {
	var companies []models.Company
	user := currentUser(c)
	var businessReq BusinessReq
	var business models.Business
	err := c.BindJSON(&businessReq)
	log.Println("artilceReq", businessReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		ret := db.Where("id = ? and user_id = ?", businessReq.ID, user.ID).Find(&business)
		if ret.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": ret.Error.Error(),
			})
		} else {
			business.Name = businessReq.Name
			business.Description = businessReq.Description
			db.Save(&business)
			fmt.Println("company_ids: ", businessReq.CompanyIds)
			db.Where("id IN ?", businessReq.CompanyIds).Find(&companies)
			db.Model(&business).Association("Companies").Append(companies)
			db.Model(&user).Association("Companies").Append(companies)
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		}
	}
}

// BusinessRespStruct list businesses
type BusinessRespStruct struct {
	ID           uint
	Name         string
	CsvCount     int64
	ArticleCount int64
}

// ListBusiness get
func ListBusiness(c *gin.Context) {
	var businesses []models.Business
	businessesResp := make([]BusinessRespStruct, 0)
	user := currentUser(c)
	// page := c.DefaultQuery("page", "1")
	// offset, error := strconv.Atoi(page)
	// if error != nil {
	// 	log.Fatal("page format error", error)
	// }
	db.Preload("Businesses").Where("user_id = ? ", user.ID).Find(&businesses) //.Offset(offset * 20).Limit(20)
	for _, business := range businesses {
		businessCount := db.Model(&business).Association("Articles").Count()
		csvCount := db.Model(&business).Association("Csvs").Count()
		businessResp := BusinessRespStruct{
			ID:           business.ID,
			Name:         business.Name,
			CsvCount:     csvCount,
			ArticleCount: businessCount,
		}
		businessesResp = append(businessesResp, businessResp)
	}
	c.JSON(http.StatusOK, gin.H{
		"businesses": businessesResp,
		"message":    "ok",
	})
}

// Business get specified business by ID
func Business(c *gin.Context) {
	var business models.Business
	user := currentUser(c)
	id := c.Param("id")
	db.Where("user_id = ? AND ID = ?", user.ID, id).Find(&business)
	c.JSON(http.StatusOK, gin.H{
		"business": business,
		"message":  "ok",
	})
}

// BusinessStat compare companies's data
type BusinessStat struct {
	Name                      string
	Code                      string
	ReportName                string
	Category                  string
	ReportDate                int
	TotalRevenue              float64
	NetProfitAfterNrgalAtsolc float64
}

// BusinessStats get specified business's stats data
func BusinessStats(c *gin.Context) {
	var business models.Business
	user := currentUser(c)
	id := c.Param("id")
	db.Preload("Companies").Where("user_id = ? AND ID = ?", user.ID, id).Find(&business)
	sql := `
		select b.code, b.name, a.report_name, a.category, a.report_date,
		a.total_revenue, a.net_profit_after_nrgal_atsolc
		from report_summaries as a, companies as b
		where a.company_code = ? and b.id = ?;
	`
	var businessesStats = make([][]BusinessStat, 0)
	for _, company := range business.Companies {
		var businessStats []BusinessStat
		db.Raw(sql, company.Code, company.ID).Scan(&businessStats)
		businessesStats = append(businessesStats, businessStats)
	}
	c.JSON(http.StatusOK, gin.H{
		"business": business,
		"stats":    businessesStats,
		"message":  "ok",
	})
}

// DeleteBusiness mark business as deleted
func DeleteBusiness(c *gin.Context) {
	var session models.Session
	var business models.Business
	token := c.GetHeader("Token")
	log.Println("token", token)
	db.Where("key = ?", token).Find(&session)
	id := c.Param("id")
	db.Where("user_id = ? AND ID = ?", session.UserID, id).Find(&business)
	db.Delete(&business)
	c.JSON(http.StatusOK, gin.H{
		"business": business,
		"message":  "ok",
	})
}

func currentUser(c *gin.Context) *models.User {
	var session models.Session
	token := c.GetHeader("Token")
	db.Preload("User").Where("key = ?", token).Find(&session)
	return session.User
}
