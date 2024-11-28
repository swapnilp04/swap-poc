package models


import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type CommentCategory struct {
	ID            					uint    `json:"id"`
	Name     								string `json:"name" validate:"nonzero"`
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}

func migrateCommentCategory() {
	fmt.Println("migrating Comment Category..")
	err := db.Driver.AutoMigrate(&CommentCategory{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func migrateCommentCategoryData() {
	categoryArr := [...]string{"Hostel", "Payment", "Teacher", "HomeWork", "Parents", "Exam"}
	for _, category := range categoryArr { 
        commentCategory := CommentCategory{Name: category}
        err := commentCategory.FindByName()
        if(err != nil) {
        	commentCategory.Create()
        }
    } 
}

func NewCommentCategory(commentCategoryData map[string]interface{}) *CommentCategory {
	commentCategory := &CommentCategory{}
	commentCategory.Assign(commentCategoryData)
	return commentCategory
}

func (cc *CommentCategory) Validate() error {
	if errs := validator.Validate(cc); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (cc *CommentCategory) Assign(commentCategoryData map[string]interface{}) {
	fmt.Printf("%+v\n", commentCategoryData)
	if name, ok := commentCategoryData["name"]; ok {
		cc.Name = name.(string)
	}
}

func (cc *CommentCategory) All() ([]CommentCategory, error) {
	var commentCategories []CommentCategory
	err := db.Driver.Find(&commentCategories).Error
	return commentCategories, err
}

func (cc *CommentCategory) Find() error {
	err := db.Driver.First(cc, "ID = ?", cc.ID).Error
	return err
}

func (cc *CommentCategory) FindByName() error {
	err := db.Driver.First(cc, "Name = ?", cc.Name).Error
	return err
}

func (cc *CommentCategory) Create() error {
	err := db.Driver.Create(cc).Error
	return err
}

func (cc *CommentCategory) Update() error {
	err := db.Driver.Save(cc).Error
	return err
}

func (cc *CommentCategory) Delete() error {
	err := db.Driver.Delete(cc).Error
	return err
}
