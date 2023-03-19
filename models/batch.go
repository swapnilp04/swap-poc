package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Batch struct {
	ID            	int    `json:"id"`
	Name     				string `json:"name"`
	Year      			int `json:"year"`
	Transactions  	[]Transaction
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

func (b *Batch) Validate() error {
	return nil
}

func (b *Batch) Assign(batchData map[string]interface{}) {
	fmt.Printf("%+v\n", batchData)
	if name, ok := batchData["name"]; ok {
		b.Name = name.(string)
	}

	if year, ok := batchData["year"]; ok {
		b.Year = int(year.(int64))
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
