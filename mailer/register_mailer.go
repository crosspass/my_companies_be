package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"github.com/my-companies-be/models"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.my-companies")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}
}

// SendActiveAccount send mail to active account by link
func SendActiveAccount(u *models.User) {
	from := viper.GetString("from")
	password := viper.GetString("password")

	// Receiver email address.
	to := []string{
		u.Email,
	}

	// smtp server configuration.
	smtpHost := viper.GetString("smtpHost")
	smtpPort := viper.GetString("smtpPort")
	smtpAddress := smtpHost + ":" + smtpPort

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("register.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Token string
	}{
		Token: u.RegisterToken,
	})

	// Sending email.
	err := smtp.SendMail(smtpAddress, auth, from, to, body.Bytes())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent!")
}
