package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Batch struct {
	ID            	uint    `json:"id"`
	Name     				string `json:"name" validate:"nonzero"`
	Year      			int `json:"year" validate:"nonzero"`
	StandardsCount  int64 `json:"standards_count"`
	IsDefault  			bool `json:"is_default" gorm:"default:false"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateBatch() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&Batch{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewBatch(batchData map[string]interface{}) *Batch {
	batch := &Batch{}
	batch.Assign(batchData)
	return batch
}

func GetDefaultBatch() (Batch, error) {
	batch := Batch{}
	err := db.Driver.Where("is_default = true").Find(&batch).Error
	return batch, err
}

func (b *Batch) Validate() error {
	if errs := validator.Validate(b); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (b *Batch) Assign(batchData map[string]interface{}) {
	fmt.Printf("%+v\n", batchData)
	if name, ok := batchData["name"]; ok {
		b.Name = name.(string)
	}

	if year, ok := batchData["year"]; ok {
		b.Year = int(year.(float64))
	}

	if isDefault, ok := batchData["is_default"]; ok {
		b.IsDefault = isDefault.(bool)
		if b.IsDefault {
			BatchDefaultFalse()
		}
	}
}

func (b *Batch) All() ([]Batch, error) {
	var batchs []Batch
	err := db.Driver.Find(&batchs).Error
	return batchs, err
}

func (b *Batch) Find() error {
	err := db.Driver.First(b, "ID = ?", b.ID).Error
	return err
}

func (b *Batch) Create() error {
	//transaction block
	err := db.Driver.Create(b).Error
	return err
}

func (b *Batch) Update() error {
	err := db.Driver.Save(b).Error
	return err
}

func (b *Batch) Delete() error {
	err := db.Driver.Delete(b).Error
	return err
}

func (b *Batch) GetBatchStandards() ([]BatchStandard, error) {
	var batchStandards []BatchStandard
	err := db.Driver.Preload("Standard").Where("batch_id = ?", b.ID).Find(&batchStandards).Error
	return batchStandards, err
}

func (b *Batch) GetBatchStandard(batchID uint) (BatchStandard, error) {
	var batchStandard BatchStandard
	err := db.Driver.Preload("Standard").First(&batchStandard, "batch_id = ? and id = ?", b.ID, batchID).Error
	return batchStandard, err
}

func BatchDefaultFalse() {
	db.Driver.Exec("Update batches set is_default = false")
}