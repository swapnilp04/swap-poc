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
	IsActive        bool `json:"is_active" gorm:"default:true"`
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
	err := db.Driver.Preload("Standard").Preload("Batch").First(bs, "ID = ?", bs.ID).Error
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

func (bs *BatchStandard) Deactivate() error {
	err := db.Driver.Model(&bs).Update("is_active", false).Error
	return err
}

func (bs *BatchStandard) Activate() error {
	err := db.Driver.Model(&bs).Update("is_active", true).Error
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

func (bs *BatchStandard) GetSubjects() ([]Subject, error) {
	subjects := []Subject{}
	err := db.Driver.Where("standard_id = ?", bs.StandardId).Find(&subjects).Error	
	return subjects, err
}

func (bs *BatchStandard) GetTeachersLogs() ([]TeacherLog, error) {
	var teacherLogs []TeacherLog
	err := db.Driver.Where("batch_standard_id = ?", bs.ID).Find(&teacherLogs).Error
	return teacherLogs, err
}

func (bs *BatchStandard) GetBatchStandardLogs(page int, searchTeacher string, searchSubject string, searchDate string) ([]TeacherLog, error) {
	var teachersLogs []TeacherLog

	query := db.Driver.Limit(10).Preload("BatchStandard.Standard").Preload("Subject").Preload("Teacher").
	Preload("LogCategory").Where("batch_standard_id = ?", bs.ID)


	if searchTeacher != "" {
		query = query.Where("teacher_id = ?", searchTeacher)
	}

	if searchSubject != "" {
		query = query.Where("subject_id = ?", searchSubject)
	}

	if searchDate != "" {
		startDate, _ := time.Parse("2/1/2006", searchDate)
		year, month, day := startDate.Date()
		endDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
		query = query.Where("log_date >= ? and log_date <= ?", startDate, endDate)
	}

	err := query.Offset((page-1) * 10).Order("log_date DESC, start_hour DESC").Find(&teachersLogs).Error
	return teachersLogs, err
}

func (bs *BatchStandard) AllBatchStandardLogsCount(searchTeacher string, searchSubject string, searchDate string) (int64, error) {
	var count int64
	query := db.Driver.Model(&TeacherLog{}).Where("batch_standard_id = ?", bs.ID)
	
	if searchTeacher != "" {
		query = query.Where("teacher_id = ?", searchTeacher)
	}

	if searchSubject != "" {
		query = query.Where("subject_id = ?", searchSubject)
	}

	if searchDate != "" {
		startDate, _ := time.Parse("2/1/2006", searchDate)
		year, month, day := startDate.Date()
		endDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
		query = query.Where("log_date >= ? and log_date <= ?", startDate, endDate)
	}

	err := query.Count(&count).Error
	return count, err
}

func (bs *BatchStandard) ReportLogs(searchDate string) ([]TeacherLog, error) {
	var teachersLogs []TeacherLog
	query := db.Driver.Preload("BatchStandard.Standard").Preload("Subject").Preload("Teacher").
	Preload("LogCategory").Where("batch_standard_id = ?", bs.ID)
	
		
	if searchDate != "" {
		startDate, _ := time.Parse("2/1/2006", searchDate)
		year, month, day := startDate.Date()
		endDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
		query = query.Where("log_date >= ? and log_date <= ?", startDate, endDate)
	}

	err := query.Order("start_hour DESC").Find(&teachersLogs).Error
	return teachersLogs, err
}

func (bs *BatchStandard) ReportMonthlyLogs(searchDate string) ([]TeacherLog, error) {
	var teachersLogs []TeacherLog
	query := db.Driver.Preload("BatchStandard.Standard").Preload("Subject").Preload("Teacher").
	Preload("LogCategory").Where("batch_standard_id = ?", bs.ID)
	
		
	if searchDate != "" {
		reportDate, _ := time.Parse("2/1/2006", searchDate)
		year, month, _ := reportDate.Date()
		startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0,1,0)
		query = query.Where("log_date >= ? and log_date <= ?", startDate, endDate)
	}

	err := query.Order("log_date ASC").Find(&teachersLogs).Error
	return teachersLogs, err
}

func (bs *BatchStandard) GetExams() ([]Exam, error) {
	var exams []Exam
	err := db.Driver.Preload("Standard").Preload("Subject").Where("batch_standard_id = ? AND exam_status != 'Created'", bs.ID).Order("id desc").Find(&exams).Error
	return exams, err
}	


