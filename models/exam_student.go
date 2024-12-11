package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type ExamStudent struct {
	ID            	uint    `json:"id"`
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
	if marks, ok := examStudentData["marks"]; ok {
		es.Marks = float32(marks.(float64))
	}

	if is_present, ok := examStudentData["is_present"]; ok {
		es.IsPresent = is_present.(bool)
		fmt.Println(es.IsPresent)
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

func (es *ExamStudent) UpdateMarks() error {
	err := db.Driver.Model(&es).Omit("Student").Omit("Exam").Updates(map[string]interface{}{"marks": es.Marks, "is_present": es.IsPresent}).Error
	return err
}

func (es *ExamStudent) Delete() error {
	err := db.Driver.Delete(es).Error
	return err
}
