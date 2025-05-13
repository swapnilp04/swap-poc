package models

import (
	"fmt"
	"swapnil-ex/models/db"
	//"swapnil-ex/swapErr"
	
	"gorm.io/gorm"
	//"github.com/pkg/errors"
	"time"
)

type ParentStudent struct {
	ID              uint    `json:"id"` 
	ParentID 				uint `json:"parent_id" validate:"nonzero"` 
	Parent 					Parent `validate:"-"`
	StudentID 			uint `json:"student_id" validate:"nonzero"` 
	Student 				Student `validate:"-"`
	Active					bool `json:"active" gorm:"default:true"` 
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateParentStudent() {
	fmt.Println("migrating ParentStudent..")
	err := db.Driver.AutoMigrate(&ParentStudent{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func (ps *ParentStudent) Find() error {
	err := db.Driver.First(ps, "ID = ?", ps.ID).Error
	return err
}

func (ps *ParentStudent) DeactiveParentStudent() error {
	err := db.Driver.Model(ps).Updates(map[string]interface{}{"active": false}).Error
	return err
}

func (ps *ParentStudent) ActiveParentStudent() error {
	err := db.Driver.Model(ps).Updates(map[string]interface{}{"active": true}).Error
	return err
}

func (ps *ParentStudent) Save() error {
	err := db.Driver.Save(ps).Error
	return err
}

func (ps *ParentStudent) All() ([]ParentStudent, error) {
	var ParentStudents []ParentStudent
	err := db.Driver.Where("parent_id = ?", ps.ParentID).Find(&ParentStudents).Error
	return ParentStudents, err
}

func (ps *ParentStudent) Delete() error {
	err := db.Driver.Delete(ps).Error
	return err
}

func (ps *ParentStudent) Load() error {
	err := db.Driver.Find(ps, "id = ?", ps.ID).Error
	return err
}

func (ps *ParentStudent) GetExamStudents(page int)  ([]ExamStudent, error) {
	var examStudents []ExamStudent
	err := db.Driver.Limit(10).Preload("Exam.Subject").Preload("Exam.Standard").Preload("Exam.Batch").Offset((page - 1) * 10).Where("student_id = ?", ps.StudentID).Omit("Student").Order("id desc").Find(&examStudents).Error
	return examStudents, err
}

func (ps *ParentStudent) GetComments(page int) ([]Comment, error) {
	comment := Comment{}
	comments, err := comment.AllByStudent(ps.StudentID, page)
	return comments, err
}

func (ps *ParentStudent) GetTransactions() ([]Transaction, error) {
	transactions := []Transaction{}
	err := db.Driver.Where("student_id = ?", ps.StudentID).Find(&transactions).Error
	return transactions, err
}

func (ps *ParentStudent) GetAccounts() ([]StudentAccount, error) {
	studentAccounts := []StudentAccount{}
	err := db.Driver.Where("student_id = ?", ps.StudentID).Find(&studentAccounts).Error
	return studentAccounts, err
}

func (ps *ParentStudent) GetStudentAccountBalance() (float64, float64) {
	studentAccounts, err := ps.Student.GetStudentAccounts()
	var totalDebits = 0.0
	var totalCredits = 0.0
	if err == nil {
		for _, studentAccount := range studentAccounts {
			if studentAccount.TransactionType == "debit" {
				totalDebits = totalDebits + studentAccount.Amount
			} else {
				totalCredits = totalCredits + studentAccount.Amount
			}
		}
	}
	return totalDebits, totalCredits
}

func (ps *ParentStudent) GetBatchStandardStudents() ([]BatchStandardStudent, error) {
	var batchStandardStudents []BatchStandardStudent
	err := db.Driver.Where("student_id = ? ", ps.StudentID).Preload("Batch").Preload("Standard").
	Find(&batchStandardStudents).Error
	return batchStandardStudents, err
}

func (ps *ParentStudent) GetStudentHostelRoommate() ([]HostelStudent, error) {
	hostelStudents, err := ps.Student.GetStudentHostelRoommate()
	return hostelStudents, err
}

func (ps *ParentStudent) GetLogAttendances(page int) ([]LogAttendance, error) {
	logAttendance := LogAttendance{}
	logAttendances, err := logAttendance.AllByStudent(ps.StudentID, page)
	return logAttendances, err
}
