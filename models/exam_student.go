package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type ExamStudent struct {
	ID            	int    `json:"id"`
	StudentID				uint `json:"student_id" validate:"nonzero"`
	Student 				Student `validate:"-"`
	ExamID					uint `json:"exam_id" validate:"nonzero"`
	Exam  					Exam `validate:"-"`
	Marks     			float32 `json:"marks"`
	Rank						int16 `json:"rank"`
	IsPresent				bool `json:"is_present" gorm:"default:true"`
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
	if errs := validator.Validate(es); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (es *ExamStudent) Assign(examStudentData map[string]interface{}) {
	fmt.Printf("%+v\n", examStudentData)
	// if name, ok := examStudentData["name"]; ok {
	// 	es.Name = name.(string)
	// }
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
