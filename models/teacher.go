package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Teacher struct {
	ID           	uint    `json:"id"`
	Name     			string `json:"name" validate:"nonzero"`
	Mobile				string `json:"mobile" validate:"nonzero,min=10,max=12"`
	AdharCard			string 	`json:"adhar_card" gorm:"adhar_card" validate:"nonzero,min=12,max=12"`
	JoiningDate		*time.Time `json:"joining_date"`
	LastDate			*time.Time `json:"last_date"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
  DeletedAt 		gorm.DeletedAt `gorm:"index"`
}

func migrateTeacher() {
	fmt.Println("migrating Teacher..")
	err := db.Driver.AutoMigrate(&Teacher{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewTeacher(teacherData map[string]interface{}) *Teacher {
	teacher := &Teacher{}
	teacher.Assign(teacherData)
	return teacher
}

func (s *Teacher) Validate() error {
	if errs := validator.Validate(s); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (t *Teacher) Assign(teacherData map[string]interface{}) {
	if name, ok := teacherData["name"]; ok {
		t.Name = name.(string)
	}

	if mobile, ok := teacherData["mobile"]; ok {
		t.Mobile = mobile.(string)
	}

	if adharCard, ok := teacherData["adhar_card"]; ok {
		t.AdharCard = adharCard.(string)
	}

	if joiningDate, ok := teacherData["joining_date"]; ok {
		s.JoiningDate, _ = time.Parse("2006-01-02T15:04:05.999999999Z", joiningDate.(string))
	}	
}

func (t *Teacher) All() ([]Teacher, error) {
	var teachers []Teacher
	err := db.Driver.Find(&teachers).Error
	return teachers, err
}

func (t *Teacher) Find() error {
	err := db.Driver.First(t, "ID = ?", t.ID).Error
	return err
}

func (t *Teacher) Create() error {
	err := db.Driver.Create(t).Error
	return err
}

func (t *Teacher) Update() error {
	err := db.Driver.Save(t).Error
	return err
}

func (t *Teacher) Delete() error {
	err := db.Driver.Delete(t).Error
	return err
}
