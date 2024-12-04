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
	BatchStandardId	int `json:"batch_standard_id"`
	BatchStandard		BatchStandard
	StandardId      uint `json:"standard_id"`
	BatchId      		uint `json:"batch_id"`
	ExamType				string `json:"exam_type"`
	ExamMarks				int `json:"exam_marks"`
	ExamTime 				int `json:"exam_time"`
	ExamDate				time.Time
	ExamStatus 			string `json:"exam_status"`
	ContactNumber  	string `json:"contact_number"`
	ExamStudents 		[]ExamStudent
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

func (e *Exam) Validate() error {
	return nil
}

func (e *Exam) Assign(examData map[string]interface{}) {
	fmt.Printf("%+v\n", examData)
	if name, ok := examData["name"]; ok {
		e.Name = name.(string)
	}
}

func (e *Exam) All() ([]Exam, error) {
	var exams []Exam
	err := db.Driver.Find(&exams).Error
	return exams, err
}

func (e *Exam) Find() error {
	err := db.Driver.First(e, "ID = ?", e.ID).Error
	return err
}

func (e *Exam) Create() error {
	err := db.Driver.Create(e).Error
	return err
}

func (e *Exam) Update() error {
	err := db.Driver.Save(e).Error
	return err
}

func (e *Exam) Delete() error {
	err := db.Driver.Delete(e).Error
	return err
}

func (e *Exam) PlotExamStudents() error {
	batchStandardStudents, err := e.BatchStandard.GetStudents()
	for _, batchStandardStudent := range batchStandardStudents {
			examStudent := &ExamStudent{StudentId: batchStandardStudent.StudentId}
			examStudent.Create()
		}
	//err := db.Driver.Find(&examStudents).Error
	return  err
}