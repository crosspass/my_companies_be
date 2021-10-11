package connect

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Db for db
var Db *gorm.DB

func init() {
	viper.SetConfigName("config")              // name of config file (without extension)
	viper.SetConfigType("yaml")                // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/my-companies/") // path to look for the config file in
	viper.AddConfigPath("$HOME/.my-companies") // call multiple times to add many search paths
	viper.AddConfigPath(".")                   // optionally look for config in the working directory
	err := viper.ReadInConfig()                // Find and read the config file
	if err != nil {                            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
  host := viper.GetString("db.host")
	dbUser := viper.GetString("db.user")
	dbName := viper.GetString("db.name")
	fmt.Println("dbUser ", dbUser, dbName)
	dsn := fmt.Sprintf("host=%s user=%s password=gorm dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",host, dbUser, dbName)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
