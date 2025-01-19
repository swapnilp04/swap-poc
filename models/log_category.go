package models


import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type LogCategory struct {
	ID            					uint    `json:"id"`
	Name     								string `json:"name" validate:"nonzero"`
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateLogCategory() {
	fmt.Println("migrating Log Category..")
	err := db.Driver.AutoMigrate(&LogCategory{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func migrateLogCategoryData() {
	categoryArr := [...]string{"Teaching", "Prectice", "Checking", "Discussion", "Exercise", "Paper", "DPP", "Exam", "Hackathon", "Revision"}
	for _, category := range categoryArr { 
        logCategory := LogCategory{Name: category}
        err := logCategory.FindByName()
        if(err != nil) {
        	logCategory.Create()
        }
    } 
}

func NewLogCategory(logCategoryData map[string]interface{}) *LogCategory {
	logCategory := &LogCategory{}
	logCategory.Assign(logCategoryData)
	return logCategory
}

func (cc *LogCategory) Validate() error {
	if errs := validator.Validate(cc); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (cc *LogCategory) Assign(logCategoryData map[string]interface{}) {
	fmt.Printf("%+v\n", logCategoryData)
	if name, ok := logCategoryData["name"]; ok {
		cc.Name = name.(string)
	}
}

func (cc *LogCategory) All() ([]LogCategory, error) {
	var logCategories []LogCategory
	err := db.Driver.Find(&logCategories).Error
	return logCategories, err
}

func (cc *LogCategory) Find() error {
	err := db.Driver.First(cc, "ID = ?", cc.ID).Error
	return err
}

func (cc *LogCategory) FindByName() error {
	err := db.Driver.First(cc, "Name = ?", cc.Name).Error
	return err
}

func (cc *LogCategory) Create() error {
	err := db.Driver.Create(cc).Error
	return err
}

func (cc *LogCategory) Update() error {
	err := db.Driver.Save(cc).Error
	return err
}

func (cc *LogCategory) Delete() error {
	err := db.Driver.Delete(cc).Error
	return err
}
