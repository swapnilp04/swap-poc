package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"swapnil-ex/swapErr"
	
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	//"github.com/pkg/errors"
	"time"
	"strings"
)

type Parent struct {
	ID              uint    `json:"id"` 
	ParentName      string 	`json:"parent_name" validate:"nonzero"`
	DisplayName			string 	`json:"display_name" validate:"nonzero"`
	Salt            string 	`json:"-"`
	Password        string 	`json:"-"`
	ConfirmPassword string 	`json:"-" gorm:"-"`
	DeviceID				string 	`json:"device_id"`
	Mobile 					string 	`json:"mobile" gorm:"mobile" validate:"nonzero,min=10,max=12"`
	Mpin   					string 	`json:"-"`
	StudentCount 		int64  	`json:"student_count"`
	Active					bool 		`json:"active" gorm:"default:true"` 
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateParent() {
	fmt.Println("migrating parent..")
	err := db.Driver.AutoMigrate(&Parent{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewParent(parentData map[string]interface{}) *Parent {
	parent := &Parent{}
	parent.Assign(parentData)
	return parent
}

func (p *Parent) Validate() error {
	if p.ConfirmPassword != p.Password {
		return swapErr.ErrPasswordMisMatch
	}
	if errs := validator.Validate(p); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (p *Parent) Assign(parentData map[string]interface{}) {
	if parentName, ok := parentData["parent_name"]; ok {
		p.ParentName = parentName.(string)
	}

	if displayName, ok := parentData["display_name"]; ok {
		p.DisplayName = displayName.(string)
	}

	if mobile, ok := parentData["mobile"]; ok {
		p.Mobile = mobile.(string)
	}
}

func (p *Parent) FindParentByParentname(parentName string) error {
	err := db.Driver.Where("parent_name = ? and active = ?", parentName, true).First(p).Error
	return err
}

func (p *Parent) ValidMpin(mpin string, deviceID string) error {
	if mpin == p.Mpin && deviceID == p.DeviceID {
		return nil	
	}
	return swapErr.ErrInvalidParent
}

func (p *Parent) Find() error {
	err := db.Driver.First(p, "ID = ?", p.ID).Error
	return err
}

func (p *Parent) Create() error {
	err := db.Driver.Create(p).Error
	return err
}

func (p *Parent) Update() error {
	err := db.Driver.Save(p).Error
	return err
}

func (p *Parent) DeactiveParent() error {
	err := db.Driver.Model(p).Updates(map[string]interface{}{"active": false}).Error
	return err
}

func (p *Parent) ActiveParent() error {
	err := db.Driver.Model(p).Updates(map[string]interface{}{"active": true}).Error
	return err
}

func (p *Parent) Save() error {
	err := db.Driver.Save(p).Error
	return err
}

func (p *Parent) All() ([]Parent, error) {
	var parents []Parent
	err := db.Driver.Find(&parents).Error
	return parents, err
}

func (p *Parent) AllWithPagination(page int, search string) ([]Parent, error) {
	var parents []Parent
	query := db.Driver.Limit(10).Offset((page - 1) * 10).Order("id desc")
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("display_name like ? or parent_name like ?", search, search)
	}
	err := query.Find(&parents).Error
	return parents, err
}

func (p *Parent) Count(search string) (int64, error) {
	var count int64
	query := db.Driver.Model(&Parent{})
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("display_name like ? or parent_name like ?", search, search)
	}
	err := query.Count(&count).Error
	return count, err
}

func (p *Parent) Delete() error {
	err := db.Driver.Delete(p).Error
	return err
}

func (p *Parent) Load() error {
	err := db.Driver.Find(p, "id = ?", p.ID).Error
	return err
}

func (p *Parent) UpdateStudentCount() error {
	var count int64
	db.Driver.Model(&ParentStudent{}).Where("parent_id = ?", p.ID).Count(&count)
	err := db.Driver.Model(&Parent{}).Where("id = ?", p.ID).Update("student_count", count).Error
	return err
}
