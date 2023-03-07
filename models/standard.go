package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Standard struct {
	ID           	int    `json:"id"`
	Name     			string `json:"name"`
	Std       		int64 `json:int64`
	CreatedAt time.Time
	UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
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
	return nil
}

func (s *Standard) Assign(standardData map[string]interface{}) {
	fmt.Printf("%+v\n", standardData)
	if name, ok := standardData["name"]; ok {
		s.Name = name.(string)
	}
}

func (s *Standard) All() ([]Standard, error) {
	var standards []Standard
	err := db.Driver.Find(&standards).Error
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
