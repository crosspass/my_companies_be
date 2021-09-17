package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
	"github.com/spf13/viper"
)

type csvRespStruct struct {
	models.Csv
	Content string
	Url     string
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
	log.Println(file.Filename, file.Size)
	csv.CompanyID = company.ID
	csv.UserID = session.UserID
	csv.OriginName = file.Filename
	csv.Name = dstName
	csv.Path = dst
	csv.Size = file.Size
	// Upload the file to specific dst.
	ctx.SaveUploadedFile(file, dst)
	// Create csv record
	db.Create(&csv)
	content, err := ioutil.ReadFile(csv.Path)
	if err != nil {
		log.Fatal("Read csv file content", err)
	}
	host := viper.GetString("host")
	url := "http://" + host + "/uploads/" + csv.Name
	csvResp := csvRespStruct{csv, string(content), url}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"csv":     csvResp,
	})
}

// IndexCsv list user's csvs for specified company
func IndexCsv(ctx *gin.Context) {
	var session models.Session
	var company models.Company
	var csvs []models.Csv
	var csvsResp []csvRespStruct
	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	db.Where("code = ?", ctx.Query("code")).Find(&company)
	db.Where("user_id = ? AND company_id = ?", session.UserID, company.ID).Find(&csvs)
	for _, csv := range csvs {
		content, err := ioutil.ReadFile(csv.Path)
		if err != nil {
			log.Fatal("Read csv file content", err)
		}
		host := viper.GetString("host")
		url := "http://" + host + "/uploads/" + csv.Name
		csvsResp = append(csvsResp, csvRespStruct{csv, string(content), url})
	}
	// single file
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"csvs":    csvsResp,
	})
}

// UpdateCSVFile update csv file
func UpdateCSVFile(ctx *gin.Context) {
	var session models.Session
	var csv models.Csv
	token := ctx.GetHeader("Token")
	db.Where("key = ?", token).Find(&session)
	db.Find(&csv, ctx.Param("id"))
	// single file
	file, _ := ctx.FormFile("file")
	timestamp := time.Now().Local().UnixNano()
	dstName := strconv.FormatInt(timestamp, 10) + "_" + file.Filename
	dst := viper.GetString("uploads") + dstName
	log.Println(file.Filename, file.Size)
	if err := os.Remove(csv.Path); err != nil {
		log.Printf("remove file %s error %s", csv.Path, err)
	}
	csv.OriginName = file.Filename
	csv.Name = dstName
	csv.Path = dst
	csv.Size = file.Size

	// Upload the file to specific dst.
	ctx.SaveUploadedFile(file, dst)
	// Create csv record
	db.Save(&csv)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
