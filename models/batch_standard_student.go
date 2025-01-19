package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type BatchStandardStudent struct {
	ID            		uint `json:"id"`
	BatchId       		uint `json:"batch_id" validate:"nonzero"`
	StandardId    		uint `json:"standard_id" validate:"nonzero"`
	StudentId					uint `json:"student_id" validate:"nonzero"`
	BatchStandardId 	uint `json:"batch_standard_id" validate:"nonzero"`
	Standard 					Standard
	Batch 						Batch
	BatchStandard     BatchStandard
	Student 					Student
	Fee 							float64 `json:"fee" validate:"nonzero"`
	CreatedAt 				time.Time
	UpdatedAt 				time.Time
  DeletedAt 				gorm.DeletedAt `gorm:"index"`
}

func migrateBatchStandardStudent() {
	fmt.Println("migrating batch standard student..")
	err := db.Driver.AutoMigrate(&BatchStandardStudent{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewBatchStandardStudent(batchStandardStudentData map[string]interface{}, student *Student) *BatchStandardStudent {
	batchStandardStudent := &BatchStandardStudent{}
	batchStandardStudent.Assign(batchStandardStudentData, student)
	return batchStandardStudent
}

func GetBatchStandardId(batchStandardStudentData map[string]interface{}) uint {
	if batch_standard_id, ok := batchStandardStudentData["batch_standard_id"]; ok {
		return uint(batch_standard_id.(float64))
	}	
	return 0
}

func (bss *BatchStandardStudent) Validate() error {
	if errs := validator.Validate(bss); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (bs *BatchStandardStudent) Assign(batchStandardStudentData map[string]interface{}, student *Student) {
	if batch_id, ok := batchStandardStudentData["batch_id"]; ok {
		bs.BatchId = uint(batch_id.(float64))
	}

	if standard_id, ok := batchStandardStudentData["standard_id"]; ok {
		bs.StandardId = uint(standard_id.(float64))
	}
	// if student_id, ok := batchStandardStudentData["student_id"]; ok {
	// 	bs.StudentId = uint(student_id.(float64))
	// }
}

func (bs *BatchStandardStudent) All(studentId uint) ([]BatchStandardStudent, error) {
	var batchStandardStudents []BatchStandardStudent
	err := db.Driver.Preload("Standard").Preload("Batch").Where("student_id = ?", studentId).Find(&batchStandardStudents).Error
	return batchStandardStudents, err
}

func (bs *BatchStandardStudent) Find() error {
	err := db.Driver.First(bs, "ID = ?", bs.ID).Error
	return err
}

func (bss *BatchStandardStudent) Create() error {
	err := db.Driver.Where(BatchStandardStudent{StudentId: bss.StudentId, BatchStandardId: bss.BatchStandardId}).
	Assign(BatchStandardStudent{StandardId: bss.StandardId, BatchId: bss.BatchId}).FirstOrCreate(bss).Error
	bss.updateCount()
	err = bss.AddTransaction()
	return err
}

func (bs *BatchStandardStudent) updateCount() {
	var count int64
	db.Driver.Model(&BatchStandardStudent{}).Where("batch_standard_id = ?", bs.BatchStandardId).Count(&count)
	db.Driver.Model(&BatchStandard{}).Where("id = ?", bs.BatchStandardId).Update("students_count", count)
}

func (bs *BatchStandardStudent) Update() error {
	err := db.Driver.Save(bs).Error
	return err
}

func (bs *BatchStandardStudent) Delete() error {
	err := db.Driver.Delete(bs).Error
	return err
}

func (bss *BatchStandardStudent) AddTransaction() error{
	transaction := &Transaction{}
	batchStandard := &BatchStandard{ID: bss.BatchStandardId}
	err := batchStandard.Find()
	if err != nil {
		return err
	}
	transactionCategory, err := batchStandard.GetTransactionCategory()
	if err != nil {
		return err
	}
	name := "New Adminission " + batchStandard.Batch.Name + "-" + batchStandard.Standard.Name
	transactionData := map[string]interface{}{"name": name, "student_id": float64(bss.StudentId), 
		"transaction_category_id": float64(transactionCategory.ID), "batch_standard_student_id": float64(bss.ID), "is_cleared": true, "transaction_type": "debit", 
		"amount": batchStandard.Fee}
	transaction.Assign(transactionData)
	err = transaction.Create()
	return err
}

func (bs *BatchStandardStudent) GetTransactions() ([]Transaction, error) {
	batchStandard := bs.BatchStandard
	transactionCategory, _ := batchStandard.GetTransactionCategory()
	transactions := []Transaction{}
	err := db.Driver.Where("TransactionCategoryId = ? AND StudentId = ? AND BatchStandardStudentId = ? AND IsCleared = ?", 
		transactionCategory.ID, bs.StandardId, bs.ID, true).Find(transactions).Error
	return transactions, err
}

func (bs *BatchStandardStudent) TotalDebits() float64 {
	transactions, err := bs.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "debit" {
				total = total + transaction.Amount
			}
		}
	}
	return total
}

func (bs *BatchStandardStudent) TotalCridits() float64 {
	transactions, err := bs.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "credit" {
				total = total + transaction.Amount
			}
		}
	}
	return total	
}
