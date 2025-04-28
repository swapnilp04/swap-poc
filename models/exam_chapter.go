package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	//"strconv"
)

type ExamChapter struct {
	ID            	uint    `json:"id"`
	ExamID					uint `json:"exam_id" validate:"nonzero"`
	Exam 						Exam  `validate:"-"`
	ChapterID				uint `json:"chapter_id" validate:"nonzero"`
	Chapter 				Chapter `validate:"-"`
	SubjectID      	uint `json:"subject_id" validate:"nonzero"`
	Subject 				Subject `validate:"-"`
	BatchStandardID uint `json:"batch_standard_id" validate:"nonzero"`
	BatchStandard 	BatchStandard `validate:"-"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateExamChapter() {
	fmt.Println("migrating Exam Chapter..")
	err := db.Driver.AutoMigrate(&ExamChapter{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewExamChapter(examChapterData map[string]interface{}) *ExamChapter {
	examChapter := &ExamChapter{}
	examChapter.Assign(examChapterData)
	return examChapter
}

func (ec *ExamChapter) Validate() error {
	if errs := validator.Validate(ec); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (ec *ExamChapter) Assign(examChapterData map[string]interface{}) {
	fmt.Printf("%+v\n", examChapterData)
	
	if batchStandardID, ok := examChapterData["batch_standard_id"]; ok {
		ec.BatchStandardID = uint(batchStandardID.(float64))
	}
	if subjectID, ok := examChapterData["subject_id"]; ok {
		ec.SubjectID = uint(subjectID.(float64))
	}

	if chapterID, ok := examChapterData["chapter_id"]; ok {
		ec.ChapterID = uint(chapterID.(float64))
	}

	if examID, ok := examChapterData["exam_id"]; ok {
		ec.ExamID = uint(examID.(float64))
	}
}

func (ec *ExamChapter) All(page int) ([]ExamChapter, error) {
	var examChapters []ExamChapter
	err := db.Driver.Preload("Subject").Preload("Batch").Preload("Chapter").Preload("Exam").Find(&examChapters).Error
	return examChapters, err
}

func (ec *ExamChapter) AllCount() (int64, error) {
	var count int64
	query := db.Driver.Model(&ExamChapter{})
	err := query.Count(&count).Error
	return count, err
}

func (ec *ExamChapter) Find() error {
	err := db.Driver.Preload("Subject").Preload("Batch").Preload("Chapter").Preload("Exam").First(ec, "ID = ?", ec.ID).Error
	return err
}

func (ec *ExamChapter) Create() error {
	err := db.Driver.Omit("Subject, Exam, Chapter, BatchStandard").Create(ec).Error
	return err
}

func (ec *ExamChapter) Update() error {
	//err := db.Driver.Updates(e).Error
	err := db.Driver.Omit("Subject, Exam, Chapter, BatchStandard").Session(&gorm.Session{FullSaveAssociations: false}).Updates(&ec).Error
	return err
}

func (ec *ExamChapter) Delete() error {
	err := db.Driver.Delete(ec).Error
	return err
}