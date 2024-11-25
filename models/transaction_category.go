package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type TransactionCategory struct {
	ID            					uint    `json:"id"`
	Name     								string `json:"name" validate:"nonzero"`
	HostelId								uint `json:"hostel_id"`
	BatchId									uint `json:"batch_id"`
	BatchStandardId         uint `json:"batch_standard_id"`
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
	
func migrateTransactionDiscountCategory() {
	_, err :=  GetDiscountTransactionCategory()
	if err != nil {
		transactionCategory := TransactionCategory{Name: "Discount"}
			db.Driver.Create(&transactionCategory)
	}
}

func NewTransactionCategory(transactionCategoryData map[string]interface{}) *TransactionCategory {
	transactionCategory := &TransactionCategory{}
	transactionCategory.Assign(transactionCategoryData)
	return transactionCategory
}

func GetDiscountTransactionCategory() (TransactionCategory, error) {
	tc := &TransactionCategory{Name: "Discount"}
	err :=  tc.FindByName()
	return *tc, err
}

func (t *TransactionCategory) Validate() error {
	if errs := validator.Validate(t); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (t *TransactionCategory) Assign(transactionCategoryData map[string]interface{}) {
	fmt.Printf("%+v\n", transactionCategoryData)
	if name, ok := transactionCategoryData["name"]; ok {
		t.Name = name.(string)
	}

	if batchId, ok := transactionCategoryData["batch_id"]; ok {
		t.BatchId = batchId.(uint)
	}

	if batchStandardId, ok := transactionCategoryData["batch_standard_id"]; ok {
		t.BatchStandardId = batchStandardId.(uint)
	}

	if hostelId, ok := transactionCategoryData["hostel_id"]; ok {
		t.HostelId = hostelId.(uint)
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

func (t *TransactionCategory) FindByName() error {
	err := db.Driver.First(t, "Name = ?", t.Name).Error
	return err
}

func (t *TransactionCategory) Create() error {
	err := db.Driver.Create(t).Error
	return err
}

func (t *TransactionCategory) Update() error {
	err := db.Driver.Save(t).Error
	return err
}

func (t *TransactionCategory) Delete() error {
	err := db.Driver.Delete(t).Error
	return err
}