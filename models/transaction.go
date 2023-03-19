package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Transaction struct {
	ID            	int    `json:"id"`
	Name     				string `json:"name"`
	HostelId				int `json:"hostel_id"`
	RoomId      		int `json:"room_id"`
	ContactNumber  	string `json:"contact_number"`
	StudentId				uint `json:"student_id"`
	BatchId				uint `json:"batch_id"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
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
}

func (t *Transaction) All() ([]Transaction, error) {
	var transactions []Transaction
	err := db.Driver.Find(&transactions).Error
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
