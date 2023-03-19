package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type ExamStudent struct {
	ID            	int    `json:"id"`
	Name     				string `json:"name"`
	HostelId				int `json:"hostel_id"`
	RoomId      		int `json:"room_id"`
	ContactNumber  	string `json:"contact_number"`
	StudentId				uint `json:"student_id"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateExamStudent() {
	fmt.Println("migrating Exam Student..")
	err := db.Driver.AutoMigrate(&ExamStudent{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewExamStudent(examStudentData map[string]interface{}) *ExamStudent {
	examStudent := &ExamStudent{}
	examStudent.Assign(examStudentData)
	return examStudent
}

func (es *ExamStudent) Validate() error {
	return nil
}

func (es *ExamStudent) Assign(examStudentData map[string]interface{}) {
	fmt.Printf("%+v\n", examStudentData)
	if name, ok := examStudentData["name"]; ok {
		es.Name = name.(string)
	}
}

func (es *ExamStudent) All() ([]ExamStudent, error) {
	var examStudents []ExamStudent
	err := db.Driver.Find(&examStudents).Error
	return examStudents, err
}

func (es *ExamStudent) Find() error {
	err := db.Driver.First(es, "ID = ?", es.ID).Error
	return err
}

func (es *ExamStudent) Create() error {
	err := db.Driver.Create(es).Error
	return err
}

func (es *ExamStudent) Update() error {
	err := db.Driver.Save(es).Error
	return err
}

func (es *ExamStudent) Delete() error {
	err := db.Driver.Delete(es).Error
	return err
}
