package db

import (
	//"gorm.io/driver/sqlite"
  "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Driver *gorm.DB


func init() {
	var err error
	 //dsn := "root:swapnilp04@tcp(eracord.c6daj9mtyykp.us-east-1.rds.amazonaws.com:3306)/eracord_development?charset=utf8mb4&parseTime=True&loc=Local"
	 //dsn := "root:Nswapnilp04p@tcp(127.0.0.1:3306)/eracord_development?charset=utf8mb4&parseTime=True&loc=Local"
	 dsn := "root:swapnilp04@tcp(127.0.0.1:3306)/eracord_development?charset=utf8mb4&parseTime=True&loc=Local"
	 Driver, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	 //Driver, err = gorm.Open(sqlite.Open("swapnil.db"), &gorm.Config{
	 	Logger: logger.Default.LogMode(logger.Info),
	 	DisableForeignKeyConstraintWhenMigrating: true,
	 })

	if err != nil {
		panic(err)
	}
}

func Close() {
}
