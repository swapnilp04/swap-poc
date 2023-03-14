package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Exam struct {
	ID            	int    `json:"id"`
	Name     				string `json:"name"`
	HostelId				int `json:"hostel_id"`
	RoomId      		int `json:"room_id"`
	ContactNumber  string `json:"contact_number"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateExam() {
	fmt.Println("migrating Exam..")
	err := db.Driver.AutoMigrate(&Exam{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewExam(examData map[string]interface{}) *Exam {
	exam := &Exam{}
	exam.Assign(examData)
	return exam
}

func (t *Exam) Validate() error {
	return nil
}

func (t *Exam) Assign(examData map[string]interface{}) {
	fmt.Printf("%+v\n", examData)
	if name, ok := examData["name"]; ok {
		t.Name = name.(string)
	}
}

func (t *Exam) All() ([]Exam, error) {
	var exams []Exam
	err := db.Driver.Find(&exams).Error
	return exams, err
}

func (t *Exam) Find() error {
	err := db.Driver.First(t, "ID = ?", t.ID).Error
	return err
}

func (t *Exam) Create() error {
	err := db.Driver.Create(t).Error
	return err
}

func (t *Exam) Update() error {
	err := db.Driver.Save(t).Error
	return err
}

func (t *Exam) Delete() error {
	err := db.Driver.Delete(t).Error
	return err
}
