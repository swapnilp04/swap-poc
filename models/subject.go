package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Subject struct {
	ID           	uint    `json:"id"`
	Name     			string `json:"name" gorm:"unique" validate:"nonzero"`
	StandardID    uint `json:"standard_id" validate:"nonzero"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
  DeletedAt 		gorm.DeletedAt `gorm:"index"`
}

func migrateSubject() {
	fmt.Println("migrating Subject..")
	err := db.Driver.AutoMigrate(&Subject{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewSubject(subjectData map[string]interface{}) *Subject {
	subject := &Subject{}
	subject.Assign(subjectData)
	return subject
}

func (s *Subject) Validate() error {
	if errs := validator.Validate(s); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (s *Subject) Assign(subjectData map[string]interface{}) {
	if name, ok := subjectData["name"]; ok {
		s.Name = name.(string)
	}
}

func (s *Subject) All() ([]Subject, error) {
	var subjects []Subject
	err := db.Driver.Find(&subjects).Error
	return subjects, err
}

func (s *Subject) Find() error {
	err := db.Driver.First(s, "ID = ?", s.ID).Error
	return err
}

func (s *Subject) Create() error {
	err := db.Driver.Create(s).Error
	return err
}

func (s *Subject) Update() error {
	err := db.Driver.Save(s).Error
	return err
}

func (s *Subject) Delete() error {
	err := db.Driver.Delete(s).Error
	return err
}
