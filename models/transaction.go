package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"swapnil-ex/swapErr"
	"gopkg.in/validator.v2"
)

type Transaction struct {
	ID            					uint    `json:"id"`
	Name     								string `json:"name" validate:"nonzero"`
	StudentId								uint `json:"student_id" validate:"nonzero"`
	HostelStudentId					uint `json:"hostel_student_id"`
	TransactionCategoryId   uint `json:"transaction_category_id"`
	BatchStandardStudentId	uint `json:"batch_standard_student_id"`
	PaidBy 									string `json:"paid_by" validate:"nonzero"`
	PaymentMode 						string `json:"payment_mode" validate:"nonzero"`
	IsCleared 							bool `json:"is_cleared" gorm:"default:false"`
	IsChecked 							bool `json:"is_checked" gorm:"default:false"`
	TransactionType         string `json:"transaction_type" gorm:"default:'debit'" validate:"nonzero"`
	Amount       						float64 `json:"amount" validate:"nonzero"`
	RecieptUrl  						string `json:"receipt_url"`
	Student 								Student
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateTransaction() {
	fmt.Println("migrating Transaction..")
	err := db.Driver.AutoMigrate(&Transaction{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewTransaction(transactionData map[string]interface{}, student Student) *Transaction {
	transaction := &Transaction{Student: student}
	transaction.Assign(transactionData)
	return transaction
}

func (t *Transaction) Validate() error {
	if errs := validator.Validate(t); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (t *Transaction) Assign(transactionData map[string]interface{}) {
	fmt.Printf("%+v\n", transactionData)
	if name, ok := transactionData["name"]; ok {
		t.Name = name.(string)
	}
	if studentId, ok := transactionData["student_id"]; ok {
		t.StudentId = uint(studentId.(float64))
	}
	if hostelStudentId, ok := transactionData["hostel_student_id"]; ok {
		t.HostelStudentId = uint(hostelStudentId.(float64))
	}
	if transactionCategoryId, ok := transactionData["transaction_category_id"]; ok {
		t.TransactionCategoryId = uint(transactionCategoryId.(float64))
	}
	if batchStandardStudentId, ok := transactionData["batch_standard_student_id"]; ok {
		t.BatchStandardStudentId = uint(batchStandardStudentId.(float64))
	}
	if isCleared, ok := transactionData["is_cleared"]; ok {
		t.IsCleared = isCleared.(bool)
	}

	if paymentMode, ok := transactionData["payment_mode"]; ok {
		t.PaymentMode = paymentMode.(string)
	}

	if paidBy, ok := transactionData["paid_by"]; ok {
		t.PaidBy = paidBy.(string)
	}

	if transactionType, ok := transactionData["transaction_type"]; ok {
		t.TransactionType = transactionType.(string)
	}

	if amount, ok := transactionData["amount"]; ok {
		t.Amount = amount.(float64)
	}
}



func (t *Transaction) AllStudents(page int, ids []uint) ([]Transaction, error) {
	var transactions []Transaction
	query := db.Driver.Preload("Student")
	if len(ids) > 0 {	
		query = query.Where("student_id in (?)", ids)
	}
	err := query.Limit(10).Offset((page - 1) * 10).Order("id desc").Find(&transactions).Error
	return transactions, err
}

func (t *Transaction) Count(ids []uint) (int64, error) {
	var count int64
	query := db.Driver.Model(&Transaction{})
	if len(ids) > 0 {	
		query = query.Where("student_id in (?)", ids)
	}
	err := query.Count(&count).Error
	return count, err
}

func (t *Transaction) All(studentId int) ([]Transaction, error) {
	var transactions []Transaction
	err := db.Driver.Where("student_id = ?", studentId).Find(&transactions).Error
	return transactions, err
}

func (t *Transaction) Find() error {
	err := db.Driver.First(t, "ID = ?", t.ID).Error
	return err
}

func (t *Transaction) Create() error {
	err := db.Driver.Create(t).Error
	return err
}

func (t *Transaction) Update() error {
	err := db.Driver.Save(t).Error
	return err
}

func (t *Transaction) Delete() error {
	err := db.Driver.Delete(t).Error
	return err
}

func (t *Transaction) CheckedTransaction() error {
		if t.IsChecked != true {
			t.IsChecked = true
			err := t.Update()
			return err
		} else {
			return swapErr.ErrAlreadyChecked
		}
}