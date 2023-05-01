package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"swapnil-ex/swapErr"
)

type Transaction struct {
	ID            					uint    `json:"id"`
	Name     								string `json:"name"`
	StudentId								uint `json:"student_id"`
	HostelStudentID					uint `json:"hostel_student_id"`
	TransactionCategoryId   uint `json:"transaction_category_id"`
	BatchStandardStudentId	uint `json:"batch_standard_student_id"`
	PaidBy 									string `json:"paid_by"`
	PaymentMode 						string `json:"payment_mode"`
	IsCleared 							bool `json:"is_cleared"`
	IsChecked 							bool `json:"is_checked"`
	TransactionType         string `json:"transaction_type" "default:'debit'"`
	Amount       						float64 `json:"amount"`
	RecieptUrl  						string `json:"receipt_url"`
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

func NewTransaction(transactionData map[string]interface{}) *Transaction {
	transaction := &Transaction{}
	transaction.Assign(transactionData)
	return transaction
}

func (t *Transaction) Validate() error {
	return nil
}

func (t *Transaction) Assign(transactionData map[string]interface{}) {
	fmt.Printf("%+v\n", transactionData)
	if name, ok := transactionData["name"]; ok {
		t.Name = name.(string)
	}
	if studentId, ok := transactionData["student_id"]; ok {
		t.StudentId = studentId.(uint)
	}
	if hostelStudentId, ok := transactionData["hostel_student_id"]; ok {
		t.HostelStudentID = hostelStudentId.(uint)
	}
	if transactionCategoryId, ok := transactionData["transaction_category_id"]; ok {
		t.TransactionCategoryId = transactionCategoryId.(uint)
	}
	if batchStandardStudentId, ok := transactionData["batch_standard_student_id"]; ok {
		t.BatchStandardStudentId = batchStandardStudentId.(uint)
	}
	if isCleared, ok := transactionData["is_cleared"]; ok {
		t.IsCleared = isCleared.(bool)
	}
	if transactionType, ok := transactionData["transaction_type"]; ok {
		t.TransactionType = transactionType.(string)
	}
	if amount, ok := transactionData["amount"]; ok {
		t.Amount = amount.(float64)
	}
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
