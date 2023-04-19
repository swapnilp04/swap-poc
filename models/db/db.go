package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Driver *gorm.DB


func init() {
	var err error
	Driver, err = gorm.Open(sqlite.Open("swapnil.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
}

func Close() {
}
