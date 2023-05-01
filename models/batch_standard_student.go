package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type BatchStandardStudent struct {
	ID            		uint `json:"id"`
	BatchId       		uint `json:"batch_id"`
	StandardId    		uint `json:"standard_id"`
	StudentId					uint `json:"student_id"`
	BatchStandardId 	uint `json:batch_standard_id""`
	Standard 					Standard
	Batch 						Batch
	BatchStandard     BatchStandard
	Student 					Student
	Fee 							float64 `json:"fee"`
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

func (bs *BatchStandardStudent) Validate() error {
	return nil
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

func (bs *BatchStandardStudent) Create() error {
	err := db.Driver.Where(BatchStandardStudent{StudentId: bs.StudentId, BatchStandardId: bs.BatchStandardId}).
	Assign(BatchStandardStudent{StandardId: bs.StandardId, BatchId: bs.BatchId}).FirstOrCreate(bs).Error
	bs.updateCount()
	//err = bs.AddTransaction()
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

func (bs *BatchStandardStudent) AddTransaction() error{
	transaction := &Transaction{}
	transactionCategory, err := bs.BatchStandard.GetTransactionCategory()
	if err == nil {
		return err
	}
	transactionData := map[string]interface{}{"Name": "New Adminission", "StudentId": bs.StudentId, 
		"TransactionCategoryId": transactionCategory.ID, "BatchStandardStudentId": bs.ID, "IsCleared": true, "TransactionType": "debit", 
		"Amount": bs.BatchStandard.Fee}
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
