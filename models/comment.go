package models


import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Comment struct {
	ID            					uint    `json:"id"`
	Comment     						string `json:"comment" validate:"nonzero"`
	StudentID								uint `json:"student_id"  validate:"nonzero"`
	HasReminder			   			bool `json:"has_reminder" gorm:"default:false"`
	ReminderOn							time.Time `json:"reminder_on"`
	CommentCategoryID				uint `json:"comment_category_id" validate:"nonzero"`
	CommentCategory 				CommentCategory
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
	if studentId, ok := commentData["student_id"]; ok {
		c.StudentID = uint(studentId.(float64))
	}
	if commentCategoryId, ok := commentData["comment_category_id"]; ok {
		c.CommentCategoryID = uint(commentCategoryId.(float64))
	}
	if hasReminder, ok := commentData["has_reminder"]; ok {
		c.HasReminder = hasReminder.(bool)
	}
	if reminderOn, ok := commentData["reminder_on"]; ok {
		c.ReminderOn, _ = time.Parse("2006-01-02T15:04:05.999999999Z", reminderOn.(string))
	}	
}

func (c *Comment) All() ([]Comment, error) {
	var comments []Comment
	err := db.Driver.Find(&comments).Error
	return comments, err
}

func (c *Comment) Find() error {
	err := db.Driver.First(c, "ID = ?", c.ID).Error
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