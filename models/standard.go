package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Standard struct {
	ID           	uint    `json:"id"`
	Name     			string `json:"name" gorm:"unique" validate:"nonzero"`
	Std       		int `json:"std" validate:"nonzero"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
  DeletedAt 		gorm.DeletedAt `gorm:"index"`
}

func migrateStandard() {
	fmt.Println("migrating Standard..")
	err := db.Driver.AutoMigrate(&Standard{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewStandard(standardData map[string]interface{}) *Standard {
	standard := &Standard{}
	standard.Assign(standardData)
	return standard
}

func (s *Standard) Validate() error {
	if errs := validator.Validate(s); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (s *Standard) Assign(standardData map[string]interface{}) {
	fmt.Printf("%+v\n", standardData)
	if name, ok := standardData["name"]; ok {
		s.Name = name.(string)
	}

	if std, ok := standardData["std"]; ok {
		s.Std = int(std.(float64))
	}
}

func (s *Standard) All() ([]Standard, error) {
	var standards []Standard
	err := db.Driver.Find(&standards).Error
	return standards, err
}

func (s *Standard) AllExcept(standardIds []uint) ([]Standard, error) {
	var standards []Standard
	var err error
	if(len(standardIds) > 0) {
		err = db.Driver.Where("id not in (?)", standardIds).Find(&standards).Error	
	} else {
		err = db.Driver.Find(&standards).Error	
	}
	return standards, err
}

func (s *Standard) Find() error {
	err := db.Driver.First(s, "ID = ?", s.ID).Error
	return err
}

func (s *Standard) Create() error {
	err := db.Driver.Create(s).Error
	return err
}

func (s *Standard) Update() error {
	err := db.Driver.Save(s).Error
	return err
}

func (s *Standard) Delete() error {
	err := db.Driver.Delete(s).Error
	return err
}
