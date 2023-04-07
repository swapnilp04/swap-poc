package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type TransactionCategory struct {
	ID            					uint    `json:"id"`
	Name     								string `json:"name"`
	HostelID								uint `json:"hostel_id"`
	BatchId									uint `json:"batch_id"`
	BatchStandardId         uint `json:"batch_standard_id"`
	RoomId      						uint `json:"room_id"`
	Transactions  					[]Transaction
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateTransactionCategory() {
	fmt.Println("migrating TransactionCategory..")
	err := db.Driver.AutoMigrate(&TransactionCategory{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewTransactionCategory(transactionCategoryData map[string]interface{}) *TransactionCategory {
	transactionCategory := &TransactionCategory{}
	transactionCategory.Assign(transactionCategoryData)
	return transactionCategory
}

func (t *TransactionCategory) Validate() error {
	return nil
}

func (t *TransactionCategory) Assign(transactionCategoryData map[string]interface{}) {
	fmt.Printf("%+v\n", transactionCategoryData)
	if name, ok := transactionCategoryData["name"]; ok {
		t.Name = name.(string)
	}
}

func (t *TransactionCategory) All() ([]TransactionCategory, error) {
	var transactionCategories []TransactionCategory
	err := db.Driver.Find(&transactionCategories).Error
	return transactionCategories, err
}

func (t *TransactionCategory) Find() error {
	err := db.Driver.First(t, "ID = ?", t.ID).Error
	return err
}

func (t *TransactionCategory) Create() error {
	err := db.Driver.Create(t).Error
	db.Commit()
	return err
}

func (t *TransactionCategory) Update() error {
	err := db.Driver.Save(t).Error
	db.Commit()
	return err
}

func (t *TransactionCategory) Delete() error {
	err := db.Driver.Delete(t).Error
	db.Commit()
	return err
}
