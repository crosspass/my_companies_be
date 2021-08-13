package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-companies-be/mailer"
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
	log.Println("user", userReq)
	log.Println("email", userReq.Email)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		user.Email = userReq.Email
		result := db.Create(&user)
		user.SetPassword(userReq.Password)
		mailer.SendActiveAccount(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    403,
				"message": result.Error.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
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

// Login user login website
func Login(c *gin.Context) {
	var userReq UserReq
	var user models.User
	err = c.BindJSON(&userReq)
	log.Println("comment", userReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
	} else {
		result := db.Preload("Session").Find(&user, "email = ?", userReq.Email)
		if result.RowsAffected == 1 {
			if user.ValidatePassword(userReq.Password) {
				token := user.GenerateToken()
				c.JSON(http.StatusOK, gin.H{
					"message": "ok",
					"token":   token,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "email or password error!",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "email or password error!",
			})
		}
	}
}
