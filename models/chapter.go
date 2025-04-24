package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Chapter struct {
	ID           	uint    `json:"id"`
	Name     			string `json:"name" validate:"nonzero"`
	StandardID    uint `json:"standard_id" validate:"nonzero"`
	SubjectID    	uint `json:"subject_id" validate:"nonzero"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
  DeletedAt 		gorm.DeletedAt `gorm:"index"`
}

func migrateChapter() {
	fmt.Println("migrating Chapter..")
	err := db.Driver.AutoMigrate(&Chapter{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewChapter(chapterData map[string]interface{}) *Chapter {
	chapter := &Chapter{}
	chapter.Assign(chapterData)
	return chapter
}

func (c *Chapter) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (c *Chapter) Assign(chapterData map[string]interface{}) {
	if name, ok := chapterData["name"]; ok {
		c.Name = name.(string)
	}
}

func (c *Chapter) All() ([]Chapter, error) {
	var chapters []Chapter
	err := db.Driver.Find(&chapters).Error
	return chapters, err
}

func (c *Chapter) Find() error {
	err := db.Driver.First(c, "ID = ?", c.ID).Error
	return err
}

func (c *Chapter) Create() error {
	err := db.Driver.Create(c).Error
	return err
}

func (c *Chapter) Update() error {
	err := db.Driver.Save(c).Error
	return err
}

func (c *Chapter) Delete() error {
	err := db.Driver.Delete(c).Error
	return err
}

func (c *Chapter) GetTeachersLogs() ([]TeacherLog, error) {
	var teacherLogs []TeacherLog
	err := db.Driver.Where("chapter_id = ?", c.ID).Find(&teacherLogs).Error
	return teacherLogs, err
}