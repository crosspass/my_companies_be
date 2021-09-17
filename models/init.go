package models

import (
	"github.com/my-companies-be/connect"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// viper.SetConfigName("config")              // name of config file (without extension)
	// viper.SetConfigType("yaml")                // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/my-comopanies/") // path to look for the config file in
	// viper.AddConfigPath("$HOME/.my-companies") // call multiple times to add many search paths
	// viper.AddConfigPath(".")                   // optionally look for config in the working directory
	// err := viper.ReadInConfig()                // Find and read the config file
	// if err != nil {                            // Handle errors reading the config file
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }
	// dbUser := viper.GetString("db.user")
	// dbName := viper.GetString("db.name")
	// fmt.Println("dbUser ", dbUser)
	// dsn := fmt.Sprintf("host=localhost user=%s password=gorm dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai", dbUser, dbName)
	// db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil { // Handle errors reading the config file
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }
	db = connect.Db
}
