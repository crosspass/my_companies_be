package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn = "host=localhost user=wu password=gorm dbname=my_companies port=5432 sslmode=disable TimeZone=Asia/Shanghai"
var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

// UserReq for model
type UserReq struct {
	UserName          string `json:"user_name"`
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

// RegisterUser for user register
// POST /users/register
func RegisterUser(c *gin.Context) {
	var userReq UserReq
	var user models.User
	err := c.BindJSON(&userReq)
	message := c.PostForm("user_name")
	log.Println("user", userReq)
	log.Println("user_name", message)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		user.Email = userReq.Email
		user.UserName = userReq.UserName
		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    403,
				"message": result.Error.Error(),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    200,
				"message": user,
			})
		}
	}
}

// ActiveUser active user
func ActiveUser(ctx *gin.Context) {
	token := ctx.Query("token")
	var user models.User
	err := db.Where("register_token = ?", token).Find(&user).Error
	user.Active(db)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": err,
	})
}
