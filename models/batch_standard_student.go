package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type BatchStandardStudent struct {
	ID            int    `json:"id"`
	BatchId       int
	Batch 				Batch
	StandardId    int
	gstandard 		Standard
	StudentId			uint `json:"student_id"`
	Student 			Student
	Fee 					float64 `json:"fee"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
  DeletedAt 		gorm.DeletedAt `gorm:"index"`
}

func migrateBatchStandardStudent() {
	fmt.Println("migrating batch standard student..")
	err := db.Driver.AutoMigrate(&BatchStandardStudent{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewBatchStandardStudent(batchStandardStudentData map[string]interface{}) *BatchStandardStudent {
	batchStandardStudent := &BatchStandardStudent{}
	batchStandardStudent.Assign(batchStandardStudentData)
	return batchStandardStudent
}

func (bs *BatchStandardStudent) Validate() error {
	return nil
}

func (bs *BatchStandardStudent) Assign(batchStandardStudentData map[string]interface{}) {
	fmt.Printf("%+v\n", batchStandardStudentData)
	if batch_id, ok := batchStandardStudentData["batch_id"]; ok {
		bs.BatchId = int(batch_id.(int64))
	}

	if standard_id, ok := batchStandardStudentData["standard_id"]; ok {
		bs.StandardId = int(standard_id.(int64))
	}

	if student_id, ok := batchStandardStudentData["student_id"]; ok {
		bs.StudentId = uint(student_id.(int64))
	}

}

func (bs *BatchStandardStudent) All() ([]BatchStandardStudent, error) {
	var batchStandardStudents []BatchStandardStudent
	err := db.Driver.Find(&batchStandardStudents).Error
	return batchStandardStudents, err
}

func (bs *BatchStandardStudent) Find() error {
	err := db.Driver.First(bs, "ID = ?", bs.ID).Error
	return err
}

func (bs *BatchStandardStudent) Create() error {
	err := db.Driver.Create(bs).Error
	return err
}

func (bs *BatchStandardStudent) Update() error {
	err := db.Driver.Save(bs).Error
	return err
}

func (bs *BatchStandardStudent) Delete() error {
	err := db.Driver.Delete(bs).Error
	return err
}
