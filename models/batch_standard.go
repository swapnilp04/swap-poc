package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type BatchStandard struct {
	ID            	int    `json:"id"`
	BatchId       	int
	Batch 					Batch
	StandardId      uint
	standard 				Standard
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

func NewBatchStandard(batchStandardData map[string]interface{}) *BatchStandard {
	batchStandard := &BatchStandard{}
	batchStandard.Assign(batchStandardData)
	return batchStandard
}

func (bs *BatchStandard) Validate() error {
	return nil
}

func (bs *BatchStandard) Assign(batchStandardData map[string]interface{}) {
	fmt.Printf("%+v\n", batchStandardData)
	if batch_id, ok := batchStandardData["batch_id"]; ok {
		bs.BatchId = int(batch_id.(int64))
	}

	if standard_id, ok := batchStandardData["standard_id"]; ok {
		bs.StandardId = int(standard_id.(int64))
	}

}

func (bs *BatchStandard) All() ([]BatchStandard, error) {
	var batchStandards []BatchStandard
	err := db.Driver.Find(&batchStandards).Error
	return batchStandards, err
}

func (bs *BatchStandard) Find() error {
	err := db.Driver.First(bs, "ID = ?", bs.ID).Error
	return err
}

func (bs *BatchStandard) Create() error {
	err := db.Driver.Create(bs).Error
	return err
}

func (bs *BatchStandard) Update() error {
	err := db.Driver.Save(bs).Error
	return err
}

func (bs *BatchStandard) Delete() error {
	err := db.Driver.Delete(bs).Error
	return err
}
