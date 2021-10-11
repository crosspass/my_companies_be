package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
	"github.com/spf13/viper"
)

type csvRespStruct struct {
	models.Csv
}

// UploadCSV upload csv file
func UploadCSV(ctx *gin.Context) {
	var session models.Session
	var company models.Company
	var csv models.Csv

	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	db.Where("code = ?", ctx.Query("code")).Find(&company)
	// single file
	file, _ := ctx.FormFile("file")
	timestamp := time.Now().Unix()
	dstName := strconv.FormatInt(timestamp, 10) + "_" + file.Filename
	dst := viper.GetString("uploads") + dstName
	csv.CompanyID = company.ID
	csv.UserID = session.UserID
	// Upload the file to specific dst.
	ctx.SaveUploadedFile(file, dst)
	// Create csv record
	db.Create(&csv)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"csv":     csv,
	})
}

// IndexCsv list user's csvs for specified company
func IndexCsv(ctx *gin.Context) {
	var session models.Session
	var company models.Company
	var csvs []models.Csv
	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	db.Where("code = ?", ctx.Query("code")).Find(&company)
	db.Where("user_id = ? AND company_id = ?", session.UserID, company.ID).Find(&csvs)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"csvs":    csvs,
	})
}

// UpdateCSVFile update csv file
func UpdateCSVFile(ctx *gin.Context) {
	var session models.Session
	var csv models.Csv
	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	db.Find(&csv, ctx.Param("id"))
	// Create csv record
	db.Save(&csv)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type csvReqStruct struct {
	Code      string
	Title     string `json:"title"`
	ChartType string `json:"chartType"`
	Data      string `json:"data"`
}

// UpdateCsv create csv form data
func UpdateCsv(ctx *gin.Context) {
	var session models.Session
	var company models.Company
	var csvReq csvReqStruct
	var csv models.Csv

	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	ctx.BindJSON(&csvReq)
	db.Where("code = ?", ctx.Query("code")).Find(&company)
	db.Find(&csv, ctx.Param("id"))
	csv.Title = csvReq.Title
	csv.ChartType = csvReq.ChartType
	csv.Data = csvReq.Data
	db.Save(&csv)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"csv":     csv,
	})
}

// CreateCsv create csv form data
func CreateCsv(ctx *gin.Context) {
	var session models.Session
	var company models.Company
	var csvReq csvReqStruct
	var csv models.Csv

	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	ctx.BindJSON(&csvReq)
	log.Println("csvReq.data", csvReq.Data)
	db.Where("code = ?", csvReq.Code).Find(&company)
	csv.CompanyID = company.ID
	csv.UserID = session.UserID
	csv.Title = csvReq.Title
	csv.ChartType = csvReq.ChartType
	csv.Data = csvReq.Data
	db.Create(&csv)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"csv":     csv,
	})
}
