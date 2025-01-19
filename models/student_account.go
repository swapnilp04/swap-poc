package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type StudentAccount struct {
	ID            					uint    `json:"id"`
	StudentId								uint `json:"student_id" validate:"nonzero"`
	TransactionType         string `json:"transaction_type" gorm:"default:'debit'" validate:"nonzero"`
	Amount       						float64 `json:"amount" validate:"nonzero"`
	Balance 								float64 `json:"balance" gorm:"default:0.0"`
	UserID									uint `json:"user_id" validate:"nonzero"`
	Student 								Student `validate:"-"`
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateStudentAccount() {
	fmt.Println("migrating StudentAccount..")
	err := db.Driver.AutoMigrate(&StudentAccount{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewStudentAccount(studentAccountData map[string]interface{}, student Student, transactionType string) *StudentAccount {
	studentAccount := &StudentAccount{Student: student, TransactionType: transactionType}
	studentAccount.Assign(studentAccountData)
	return studentAccount
}

func (sa *StudentAccount) Validate() error {
	if errs := validator.Validate(sa); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (sa *StudentAccount) Assign(studentAccountData map[string]interface{}) {
	if studentId, ok := studentAccountData["student_id"]; ok {
		sa.StudentId = uint(studentId.(float64))
	}
	
	if amount, ok := studentAccountData["amount"]; ok {
		sa.Amount = amount.(float64)
	}
}

func (sa *StudentAccount) AllStudentAccounts(page int, ids []uint) ([]StudentAccount, error) {
	var studentAccounts []StudentAccount
	query := db.Driver.Preload("Student")
	if len(ids) > 0 {	
		query = query.Where("student_id in (?)", ids)
	}
	err := query.Limit(10).Offset((page - 1) * 10).Order("id desc").Find(&studentAccounts).Error
	return studentAccounts, err
}

func (sa *StudentAccount) Count(ids []uint) (int64, error) {
	var count int64
	query := db.Driver.Model(&StudentAccount{})
	if len(ids) > 0 {	
		query = query.Where("student_id in (?)", ids)
	}
	err := query.Count(&count).Error
	return count, err
}

func (sa *StudentAccount) All(studentId int) ([]StudentAccount, error) {
	var studentAccounts []StudentAccount
	err := db.Driver.Where("student_id = ?", studentId).Order("id desc").Find(&studentAccounts).Error
	return studentAccounts, err
}

func (sa *StudentAccount) Find() error {
	err := db.Driver.First(sa, "ID = ?", sa.ID).Error
	return err
}

func (sa *StudentAccount) Create() error {
	err := db.Driver.Omit("Student").Create(sa).Error
	return err
}

func (sa *StudentAccount) Update() error {
	err := db.Driver.Save(sa).Error
	return err
}

func (sa *StudentAccount) Delete() error {
	err := db.Driver.Delete(sa).Error
	return err
}