package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type LogAttendance struct {
	ID           						uint `json:"id"`
	IsPresent 							bool `json:"is_present" gorm:"default:true"`
	TeacherLogID						uint `json:"teacher_log_id" validate:"nonzero"`
	StudentID 							uint `json:"student_id" validate:"nonzero"`
	Student 								Student `validate:"-"`
	TeacherLog 							TeacherLog	`validate:"-"`
	BatchStandardStudentID 	uint `json:"batch_standard_student_id" validate:"nonzero"`
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateLogAttendance() {
	fmt.Println("migrating LogAttendance..")
	err := db.Driver.AutoMigrate(&LogAttendance{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewLogAttendance(logAttendanceData map[string]interface{}) *LogAttendance {
	logAttendance := &LogAttendance{}
	logAttendance.Assign(logAttendanceData)
	return logAttendance
}

func (la *LogAttendance) Validate() error {
	if errs := validator.Validate(la); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (la *LogAttendance) Assign(logAttendanceData map[string]interface{}) {
	fmt.Printf("%+v\n", logAttendanceData)
	
	if batchStandardStudentID, ok := logAttendanceData["batch_standard_student_id"]; ok {
		la.BatchStandardStudentID = uint(batchStandardStudentID.(float64))
	}
	if studentID, ok := logAttendanceData["student_id"]; ok {
		la.StudentID = uint(studentID.(float64))
	}

	if teacherLogID, ok := logAttendanceData["teacher_log_id"]; ok {
		la.TeacherLogID = uint(teacherLogID.(float64))
	}
}

func (la *LogAttendance) All() ([]LogAttendance, error) {
	var logAttendances []LogAttendance
	err := db.Driver.Find(&logAttendances).Error
	return logAttendances, err
}

func (la *LogAttendance) Find() error {
	err := db.Driver.Preload("TeacherLog").First(la, "ID = ?", la.ID).Error
	return err
}

func (la *LogAttendance) Create() error {
	err := db.Driver.Omit("Student").Omit("TeacherLog").Create(la).Error
	return err
}

func (la *LogAttendance) Update() error {
	err := db.Driver.Omit("Student").Omit("TeacherLog").Save(la).Error
	return err
}

func (la *LogAttendance) Delete() error {
	err := db.Driver.Omit("Student").Omit("TeacherLog").Delete(la).Error
	return err
}

func (la *LogAttendance) ToggleAttendance() error {
	err := db.Driver.Omit("Student").Omit("TeacherLog").Model(&la).Update("is_present", !la.IsPresent).Error
	return err
}

func (la *LogAttendance) AllByStudent(studentId uint, page int) ([]LogAttendance, error) {
	var logAttendances []LogAttendance
	err := db.Driver.Limit(10).Preload("TeacherLog").Preload("TeacherLog.BatchStandard.Standard").Preload("TeacherLog.BatchStandard.Batch").
				Preload("TeacherLog.Subject").Preload("TeacherLog.Teacher").Preload("TeacherLog.LogCategory").Preload("TeacherLog.Chapter").Offset((page-1) * 10).
				Where("student_id = ?", studentId).Order("id desc").Find(&logAttendances).Error
	return logAttendances, err
}

func (la *LogAttendance) AllByStudentCount(studentId uint) (int64, error) {
	var count int64
	err := db.Driver.Model(&LogAttendance{}).Where("student_id = ?", studentId).Count(&count).Error
	return count, err
}