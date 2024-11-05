package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type BatchStandard struct {
	ID            	uint    `json:"id"`
	BatchId       	uint `json:"batch_id"`
	Batch 					Batch
	StandardId      uint `json:"standard_id"`
	Standard 				Standard
	Fee							float64 `json:"fee"  validate:"nonzero"`
	StudentsCount 	int64 `json:"students_count"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateBatchStandard() {
	fmt.Println("migrating batch standard..")
	err := db.Driver.AutoMigrate(&BatchStandard{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewBatchStandard(batchStandardData map[string]interface{}, batch *Batch) *BatchStandard {
	batchStandard := &BatchStandard{}
	if standard_id, ok := batchStandardData["standard_id"]; ok {
		StandardId := uint(standard_id.(float64))
		db.Driver.Where("batch_id = ? and standard_id = ?", batch.ID, StandardId).FirstOrInit(&batchStandard)
		batchStandard.StandardId = StandardId
	}

	batchStandard.Assign(batchStandardData)
	batchStandard.BatchId = batch.ID
	return batchStandard
}

func (bs *BatchStandard) Validate() error {
	if errs := validator.Validate(bs); errs != nil {
	return errs
	} else {
	return nil
	}
}

func (bs *BatchStandard) Assign(batchStandardData map[string]interface{}) {
	if fee, ok := batchStandardData["fee"]; ok {
		bs.Fee = fee.(float64)
	}
}

func (bs *BatchStandard) All(batchId uint) ([]BatchStandard, error) {
	var batchStandards []BatchStandard
	err := db.Driver.Where("batch_id = ?", batchId).Preload("Standard").Find(&batchStandards).Error
	return batchStandards, err
}

func (bs *BatchStandard) AllIds(batchId uint) ([]uint, error) {
	//var batchStandards []BatchStandard
	var ids []uint
	db.Driver.Where("batch_id = ?", batchId).Model(&BatchStandard{}).Pluck("StandardId", &ids)
	return ids, nil
}

func (bs *BatchStandard) Find() error {
	err := db.Driver.First(bs, "ID = ?", bs.ID).Error
	return err
}

func (bs *BatchStandard) Create() error {
	err := db.Driver.Save(bs).Error
	bs.updateCount()
	if err != nil {
		return err
	} else {
		err = bs.createTransactionCategory()
	}
	return err
}

func (bs *BatchStandard) updateCount() {
	var count int64
	db.Driver.Model(&BatchStandard{}).Where("standard_id = ?", bs.StandardId).Count(&count)
	db.Driver.Model(&Batch{}).Where("id = ?", bs.BatchId).Update("standards_count", count)
}

func (bs *BatchStandard) Update() error {
	err := db.Driver.Save(bs).Error
	return err
}

func (bs *BatchStandard) Delete() error {
	err := db.Driver.Delete(bs).Error
	return err
}

func (bs *BatchStandard) HasFeeAssigned() bool {
	return bs.Fee > 0.0
}

func (bs *BatchStandard) createTransactionCategory() error {
	var transactionCatetoryData = map[string]interface{}{"name": "BatchStandard", "batch_id": bs.BatchId, "batch_standard_id": bs.ID}
	transactionCategory := NewTransactionCategory(transactionCatetoryData)
	err := transactionCategory.Create()
	return err
}

func (bs *BatchStandard) GetTransactionCategory() (*TransactionCategory, error) {
	tc := &TransactionCategory{}
	err := db.Driver.Where("name like ? and batch_id = ? and batch_standard_id = ?", "BatchStandard", bs.BatchId, bs.ID).First(tc).Error
	return tc, err
}

func (bs *BatchStandard) GetStudents() ([]BatchStandardStudent, error) {
	batchStandardStudents := []BatchStandardStudent{}
	err := db.Driver.Where("batch_standard_id = ?", bs.ID).Preload("Student").Find(&batchStandardStudents).Error	
	return batchStandardStudents, err
}