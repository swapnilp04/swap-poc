package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type TeacherLog struct {
	ID           		uint `json:"id"`
	StartHour     	int `json:"start_hour" validate:"nonzero"`
	StartMinuit   	int `json:"start_minuit" validate:"nonzero"`
	EndHour       	int `json:"end_hour" validate:"nonzero"`
	EndMinuit     	int `json:"end_minuit" validate:"nonzero"`
	TeacherID     	uint `json:"teacher_id" validate:"nonzero"`
	Teacher 				Teacher
	SubjectID  			uint `json:"subject_id" validate:"nonzero"`
	Subject 				Subject
	BatchStandardID uint `json:"batch_standard_id" validate:"nonzero"`
	BatchStandard 	BatchStandard `validate:"-"`
	Comment 				string `json:"comment"`
	LogCategoryID 	uint `json:"log_category_id" validate:"nonzero"`
	LogCategory 		LogCategory
	ApprovedOn			*time.Time `json:"approved_on"`
	ApprovedBy			uint `json:"approved_by"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateTeacherLog() {
	fmt.Println("migrating TeacherLog..")
	err := db.Driver.AutoMigrate(&TeacherLog{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewTeacherLog(teachersLogData map[string]interface{}) *TeacherLog {
	teachersLog := &TeacherLog{}
	teachersLog.Assign(teachersLogData)
	return teachersLog
}

func (tl *TeacherLog) Validate() error {
	if errs := validator.Validate(tl); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (tl *TeacherLog) Assign(teachersLogData map[string]interface{}) {
	fmt.Printf("%+v\n", teachersLogData)
	if startHour, ok := teachersLogData["start_hour"]; ok {
		tl.StartHour = int(startHour.(int64))
	}

	if startMinuit, ok := teachersLogData["start_minuit"]; ok {
		tl.StartMinuit = int(startMinuit.(int64))
	}

	if endHour, ok := teachersLogData["end_hour"]; ok {
		tl.EndHour = int(endHour.(int64))
	}

	if endMinuit, ok := teachersLogData["end_minuit"]; ok {
		tl.EndMinuit = int(endMinuit.(int64))
	}

	if subjectID, ok := teachersLogData["subject_id"]; ok {
		tl.SubjectID = uint(subjectID.(float64))
	}

	if batchStandardID, ok := teachersLogData["batch_standard_id"]; ok {
		tl.BatchStandardID = uint(batchStandardID.(float64))
	}

	if logCategoryID, ok := teachersLogData["log_category_id"]; ok {
		tl.LogCategoryID = uint(logCategoryID.(float64))
	}

	if comment, ok := teachersLogData["comment"]; ok {
		tl.Comment = comment.(string)
	}
}

func (tl *TeacherLog) All() ([]TeacherLog, error) {
	var teachersLogs []TeacherLog
	err := db.Driver.Find(&teachersLogs).Error
	return teachersLogs, err
}

func (tl *TeacherLog) Find() error {
	err := db.Driver.Preload("BatchStandard, Subject, Teacher, LogCategory").First(tl, "ID = ?", tl.ID).Error
	return err
}

func (tl *TeacherLog) Create() error {
	err := db.Driver.Create(tl).Error
	return err
}

func (tl *TeacherLog) Update() error {
	err := db.Driver.Omit("BatchStandard, Subject, Teacher, LogCategory").Save(tl).Error
	return err
}

func (tl *TeacherLog) Delete() error {
	err := db.Driver.Delete(tl).Error
	return err
}