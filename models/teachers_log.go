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
	LogDate					*time.Time `json:"log_date" validate:"nonzero"`
	StartHour     	int `json:"start_hour" validate:"max=24, min=0"`
	StartMinuit   	int `json:"start_minuit" validate:"max=60,min=0"`
	EndHour       	int `json:"end_hour" validate:"max=24,min=0"`
	EndMinuit     	int `json:"end_minuit" validate:"max=60,min=0"`
	Duration       	int `json:"duration"`
	TeacherID     	uint `json:"teacher_id" validate:"nonzero"`
	Teacher 				Teacher `validate:"-"`
	SubjectID  			uint `json:"subject_id" validate:"nonzero"`
	Subject 				Subject `validate:"-"`
	ChapterID  			uint `json:"chapter_id" validate:"nonzero"`
	Chapter 				Chapter `validate:"-"`
	BatchStandardID uint `json:"batch_standard_id" validate:"nonzero"`
	BatchStandard 	BatchStandard `validate:"-"`
	Comment 				string `json:"comment"`
	LogCategoryID 	uint `json:"log_category_id" validate:"nonzero"`
	LogCategory 		LogCategory `validate:"-"`
	ApprovedOn			*time.Time `json:"approved_on"`
	ApprovedBy			uint `json:"approved_by"`
	UserID  				uint `json:"user_id" validate:"nonzero"`
	HasCombinedClass bool `json:"has_combined_class" gorm:"default:false"`
	StudentsCount 	int64 `json:"students_count"`
	AbsentCount 		int64 `json:"absent_count"`
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
	if logDate, ok := teachersLogData["log_date"]; ok {
		var time, _ = time.Parse("2006-01-02T15:04:05.999999999Z", logDate.(string))
		tl.LogDate = &time
	}

	if startHour, ok := teachersLogData["start_hour"]; ok {
		tl.StartHour = int(startHour.(float64))
	}

	if startMinuit, ok := teachersLogData["start_minuit"]; ok {
		tl.StartMinuit = int(startMinuit.(float64))
	}

	if endHour, ok := teachersLogData["end_hour"]; ok {
		tl.EndHour = int(endHour.(float64))
	}

	if endMinuit, ok := teachersLogData["end_minuit"]; ok {
		tl.EndMinuit = int(endMinuit.(float64))
	}

	if subjectID, ok := teachersLogData["subject_id"]; ok {
		tl.SubjectID = uint(subjectID.(float64))
	}

	if chapterID, ok := teachersLogData["chapter_id"]; ok {
		tl.ChapterID = uint(chapterID.(float64))
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

func (tl *TeacherLog) AssignUpdate(teachersLogData map[string]interface{}) {

	if startHour, ok := teachersLogData["start_hour"]; ok {
		tl.StartHour = int(startHour.(float64))
	}

	if startMinuit, ok := teachersLogData["start_minuit"]; ok {
		tl.StartMinuit = int(startMinuit.(float64))
	}

	if endHour, ok := teachersLogData["end_hour"]; ok {
		tl.EndHour = int(endHour.(float64))
	}

	if endMinuit, ok := teachersLogData["end_minuit"]; ok {
		tl.EndMinuit = int(endMinuit.(float64))
	}

	if subjectID, ok := teachersLogData["subject_id"]; ok {
		tl.SubjectID = uint(subjectID.(float64))
	}

	if chapterID, ok := teachersLogData["chapter_id"]; ok {
		tl.ChapterID = uint(chapterID.(float64))
	}

	if logCategoryID, ok := teachersLogData["log_category_id"]; ok {
		tl.LogCategoryID = uint(logCategoryID.(float64))
	}

	if comment, ok := teachersLogData["comment"]; ok {
		tl.Comment = comment.(string)
	}
}

func (tl *TeacherLog) All(page int, searchBatchStandard string, searchSubject string, searchTeacher string, searchDate string) ([]TeacherLog, error) {
	var teachersLogs []TeacherLog
	query := db.Driver.Limit(10).Preload("BatchStandard.Standard").Preload("BatchStandard.Batch").Preload("Subject").Preload("Teacher").
	Preload("LogCategory").Preload("Chapter")
	if searchTeacher != "" {
		query = query.Where("teacher_id = ?", searchTeacher)
	}

	if searchBatchStandard != "" {
		query = query.Where("batch_standard_id = ?", searchBatchStandard)
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

func (tl *TeacherLog) AllCount(searchBatchStandard string, searchSubject string, searchTeacher string, searchDate string) (int64, error) {
	var count int64
	query := db.Driver.Model(&TeacherLog{})
	if searchTeacher != "" {
		query = query.Where("teacher_id = ?", searchTeacher)
	}

	if searchBatchStandard != "" {
		query = query.Where("batch_standard_id = ?", searchBatchStandard)
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

func (tl *TeacherLog) Find() error {
	err := db.Driver.Preload("BatchStandard.Standard").Preload("BatchStandard.Batch").Preload("Subject").
					Preload("Teacher").Preload("LogCategory").Preload("Chapter").First(tl, "ID = ?", tl.ID).Error
	return err
}

func (tl *TeacherLog) Create() error {
	err := db.Driver.Omit("BatchStandard, Subject, Teacher, LogCategory").Create(tl).Error
	if err == nil {
		err = tl.PlotLogAttendance()
		if err == nil {
			tl.updateStudentsCount()
		}
		err = tl.CalculateDuration()
	}
	return err
}

func (tl *TeacherLog) Update() error {
	err := db.Driver.Omit("BatchStandard, Subject, Teacher, LogCategory").Save(tl).Error
	if err == nil {
		err = tl.CalculateDuration()
	}
	return err
}

func (tl *TeacherLog) Delete() error {
	err := tl.DeleteLogAttendance()
	if err == nil {
		err = db.Driver.Unscoped().Delete(tl).Error
	}
	return err
}

func(tl *TeacherLog) GetBatchStandard() (BatchStandard, error) {
	batchStandard := BatchStandard{}
	
	err := db.Driver.First(&batchStandard, "ID = ?", tl.BatchStandardID).Error
	return batchStandard, err
}

func (tl *TeacherLog) PlotLogAttendance() error {
	batchStandard, err := tl.GetBatchStandard()
	if(err != nil) {
		return err
	}

	batchStandardStudents, err := batchStandard.GetStudents()
	for _, batchStandardStudent := range batchStandardStudents {
			logAttendance := &LogAttendance{StudentID: batchStandardStudent.StudentId, TeacherLogID: tl.ID, 
												BatchStandardStudentID: batchStandardStudent.ID}
			logAttendance.Create()
		}
	return  err
}

func (tl *TeacherLog) DeleteLogAttendance() error {
	var logAttendances []LogAttendance
	err := db.Driver.Unscoped().Where("teacher_log_id = ?", tl.ID).Delete(&logAttendances).Error
	return err
}

func (tl *TeacherLog) GetLogAttendances() ([]LogAttendance, error) {
	var logAttendances []LogAttendance
	err := db.Driver.Preload("Student").Where("teacher_log_id = ?", tl.ID).Find(&logAttendances).Error
	return logAttendances, err
}

func (tl *TeacherLog) updateStudentsCount() {
	var count int64
	db.Driver.Model(&LogAttendance{}).Where("teacher_log_id = ?", tl.ID).Count(&count)
	db.Driver.Model(&TeacherLog{}).Where("id = ?", tl.ID).Update("students_count", count)
}

func (tl *TeacherLog) UpdateAbsentsCount() error{
	var count int64
	db.Driver.Model(&LogAttendance{}).Where("teacher_log_id = ? && is_present = ?", tl.ID, false).Count(&count)
	err := db.Driver.Model(&TeacherLog{}).Where("id = ?", tl.ID).Update("absent_count", count).Error
	return err
}

func (tl *TeacherLog) CreateCombinedClasses(combinedClasses []interface {}) error {
	for _, combinedClasses := range combinedClasses {
		classObj := combinedClasses.(map[string]interface{})
		newBatchStandardID := classObj["batch_standard_id"].(float64)
		newSubjectID := classObj["subject_id"].(float64)
		tLog := &TeacherLog{BatchStandardID: uint(newBatchStandardID), SubjectID: uint(newSubjectID), LogDate: tl.LogDate,
			StartHour: tl.StartHour, StartMinuit: tl.StartMinuit, EndHour: tl.EndHour, EndMinuit: tl.EndMinuit,
		 	TeacherID: tl.TeacherID, Comment: tl.Comment, LogCategoryID: tl.LogCategoryID, UserID: tl.UserID, HasCombinedClass: true}
		 err := tLog.Create()
		 if err != nil {
		 	return err
		 }
	}
	if(len(combinedClasses) > 0) { 
		err := db.Driver.Model(&tl).Updates(map[string]interface{}{"has_combined_class": true}).Error
		return err
	}
	return nil
}

func (tl *TeacherLog) CalculateDuration() error { 
	duration := 0
	hours := tl.EndHour - tl.StartHour
	hoursMin := hours * 60 // Convert into minuit 

	if(hours >= 1 && tl.StartMinuit > tl.EndMinuit) {
		duration = hoursMin - (tl.StartMinuit - tl.EndMinuit)
	} else if (hours >= 1 && tl.EndMinuit > tl.StartMinuit){
		duration = hoursMin + (tl.EndMinuit - tl.StartMinuit)
	} else {
		duration = hoursMin + (tl.EndMinuit - tl.StartMinuit)
	}
	err := db.Driver.Model(&TeacherLog{}).Where("id = ?", tl.ID).Update("duration", duration).Error
	return err
}