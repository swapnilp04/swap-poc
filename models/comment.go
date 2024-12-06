package models


import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	//"database/sql"
)

type Comment struct {
	ID            					uint    `json:"id"`
	Comment     						string `json:"comment" validate:"nonzero"`
	StudentID								uint `json:"student_id"  validate:"nonzero"`
	HasReminder			   			bool `json:"has_reminder" gorm:"default:false"`
	ReminderOn							*time.Time `json:"reminder_on"`
	CommentCategoryID				uint `json:"comment_category_id" validate:"nonzero"`
	CommentCategory 				CommentCategory `validate:"-"`
	UserID									uint `json:"user_id" validate:"nonzero"`
	User  									User `validate:"-"`
	Student  								Student `validate:"-"`
	Completed  							bool `json:"completed" gorm:"default:false"`
	CompletedOn							*time.Time `json:"completed_on"`
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateComment() {
	fmt.Println("migrating Comment..")
	err := db.Driver.AutoMigrate(&Comment{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewComment(commentData map[string]interface{}) *Comment {
	comment := &Comment{}
	comment.Assign(commentData)
	return comment
}

func (c *Comment) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (c *Comment) Assign(commentData map[string]interface{}) {
	fmt.Printf("%+v\n", commentData)
	if comment, ok := commentData["comment"]; ok {
		c.Comment = comment.(string)
	}
	if commentCategoryId, ok := commentData["comment_category_id"]; ok {
		c.CommentCategoryID = uint(commentCategoryId.(float64))
	}
	if hasReminder, ok := commentData["has_reminder"]; ok {
		c.HasReminder = hasReminder.(bool)
	}
	if reminderOn, ok := commentData["reminder_on"]; ok {
		var time, _ = time.Parse("2006-01-02T15:04:05.999999999Z", reminderOn.(string))
		c.ReminderOn = &time
	}
}

func (c *Comment) All(page int, ids []uint) ([]Comment, error) {
	var comments []Comment
	query := db.Driver.Limit(10).Preload("User").Preload("CommentCategory").Preload("Student")
	if len(ids) > 0 {	
		query = query.Where("student_id in (?)", ids)
	} else {
		query = query.Where("student_id in (?)", 0)
	}
	err := query.Offset((page-1) * 10).Order("id desc").Find(&comments).Error
	return comments, err
}

func (c *Comment) AllCount(ids []uint) (int64, error) {
	var count int64
	query := db.Driver.Model(&Comment{})
	if len(ids) > 0 {	
		query = query.Where("student_id in (?)", ids)
	}
	err := query.Count(&count).Error

	return count, err
}

func (c *Comment) MakeCompleted() error { 
		currentTime := time.Now()
		c.CompletedOn = &currentTime
		c.Completed = true
		err := c.Update()
		return err
}

func (c *Comment) AllByStudent(studentId uint, page int) ([]Comment, error) {
	var comments []Comment
	err := db.Driver.Limit(10).Preload("User").Preload("CommentCategory").Offset((page-1) * 10).Where("student_id = ?", studentId).Order("id desc").Find(&comments).Error
	return comments, err
}

func (c *Comment) AllByStudentCount(studentId uint) (int64, error) {
	var count int64
	err := db.Driver.Model(&Comment{}).Where("student_id = ?", studentId).Count(&count).Error
	return count, err
}

func (c *Comment) UpcommingComments() ([]Comment, error) {
	var comments []Comment
	time := time.Now().AddDate(0, 0, -2)
	err := db.Driver.Preload("User").Preload("CommentCategory").Preload("Student").Where("has_reminder = ? AND reminder_on > ?", true, time).Order("reminder_on asc").Find(&comments).Error
	return comments, err
}


func (c *Comment) Find() error {
	err := db.Driver.Preload("User").Preload("CommentCategory").First(c, "ID = ?", c.ID).Error
	return err
}

func (c *Comment) Create() error {
	err := db.Driver.Create(c).Error
	return err
}

func (c *Comment) Update() error {
	err := db.Driver.Save(c).Error
	return err
}

func (c *Comment) Delete() error {
	err := db.Driver.Delete(c).Error
	return err
}
