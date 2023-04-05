package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Driver, tx *gorm.DB


func init() {
	var err error
	tx, err = gorm.Open(sqlite.Open("swapnil.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	Driver = tx.Begin()
}

func Commit() error{
	return Driver.Commit().Error
}


func Close() {
}
