package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
)

type Teacher struct {
	ID           	uint    `json:"id"`
	Name     			string `json:"name" validate:"nonzero"`
	Mobile				string `json:"mobile" validate:"nonzero,min=10,max=12"`
	AdharCard			string 	`json:"adhar_card" gorm:"adhar_card" validate:"nonzero,min=12,max=12"`
	JoiningDate		*time.Time `json:"joining_date"`
	LastDate			*time.Time `json:"last_date"`
	UserID 				uint `json:"user_id"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
  DeletedAt 		gorm.DeletedAt `gorm:"index"`
}

func migrateTeacher() {
	fmt.Println("migrating Teacher..")
	err := db.Driver.AutoMigrate(&Teacher{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewTeacher(teacherData map[string]interface{}) *Teacher {
	teacher := &Teacher{}
	teacher.Assign(teacherData)
	return teacher
}

func (s *Teacher) Validate() error {
	if errs := validator.Validate(s); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (t *Teacher) Assign(teacherData map[string]interface{}) {
	if name, ok := teacherData["name"]; ok {
		t.Name = name.(string)
	}

	if mobile, ok := teacherData["mobile"]; ok {
		t.Mobile = mobile.(string)
	}

	if adharCard, ok := teacherData["adhar_card"]; ok {
		t.AdharCard = adharCard.(string)
	}

	if joiningDate, ok := teacherData["joining_date"]; ok {
		var time, _ = time.Parse("2006-01-02T15:04:05.999999999Z", joiningDate.(string))
		t.JoiningDate = &time
	}	
}

func (t *Teacher) All() ([]Teacher, error) {
	var teachers []Teacher
	err := db.Driver.Find(&teachers).Error
	return teachers, err
}

func (t *Teacher) Find() error {
	err := db.Driver.First(t, "ID = ?", t.ID).Error
	return err
}

func (t *Teacher) FindByUser() error {
	err := db.Driver.First(t, "user_id = ?", t.UserID).Error
	return err
}

func (t *Teacher) Create() error {
	err := db.Driver.Create(t).Error
	return err
}

func (t *Teacher) Update() error {
	err := db.Driver.Save(t).Error
	return err
}

func (t *Teacher) Delete() error {
	err := db.Driver.Delete(t).Error
	return err
}

func (t *Teacher) CreateUser() error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("rand.Read(salt)", err)
		return err
	}

	hash, err := scrypt.Key([]byte(t.Mobile), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return err
	}

	confirmHash, err := scrypt.Key([]byte(t.Mobile), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return err
	}

	var user User
	user.Username = t.Mobile
	user.DisplayName = t.Name
	user.Password = hex.EncodeToString(hash)
	user.ConfirmPassword = hex.EncodeToString(confirmHash)
	user.Role = "Teacher"
	user.Salt = hex.EncodeToString(salt)
	if err := user.Validate(); err != nil {
		return err
	}
	user.Save()
	db.Driver.Model(&t).Updates(Teacher{UserID: uint(user.ID)})	
	return err
}


func (t *Teacher) GetTeachersLogs(page int, searchBatchStandard string, searchSubject string, searchDate string) ([]TeacherLog, error) {
	var teachersLogs []TeacherLog

	query := db.Driver.Limit(10).Preload("BatchStandard.Standard").Preload("BatchStandard.Batch").Preload("Subject").Preload("Teacher").
	Preload("LogCategory").Where("teacher_id = ?", t.ID)


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

func (t *Teacher) AllTeachersLogsCount(searchBatchStandard string, searchSubject string, searchDate string) (int64, error) {
	var count int64
	query := db.Driver.Model(&TeacherLog{}).Where("teacher_id = ?", t.ID)
	
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

