package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type HostelStudent struct {
	ID            	uint    `json:"id"`
	Name     				string `json:"name"`
	HostelId				uint `json:"hostel_id"  validate:"nonzero"`
	HostelRoomId    uint `json:"hostel_room_id"  validate:"nonzero"`
	ContactNumber  	string `json:"contact_number"  validate:"nonzero"`
	StudentId				uint `json:"student_id"  validate:"nonzero"`
	FeeIncluded  		bool `json:"fee_included" gorm:"default:false"`
	Hostel 					Hostel
	HostelRoom      HostelRoom
	Student  				Student
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateHostelStudent() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&HostelStudent{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewHostelStudent(hostelStudentData map[string]interface{}) *HostelStudent {
	hostelStudent := &HostelStudent{}
	hostelStudent.Assign(hostelStudentData)
	return hostelStudent
}

func (hs *HostelStudent) Validate() error {
	if errs := validator.Validate(hs); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (hs *HostelStudent) Assign(hostelStudentData map[string]interface{}) {
	fmt.Printf("%+v\n", hostelStudentData)
	if name, ok := hostelStudentData["name"]; ok {
		hs.Name = name.(string)
	}

	if hostel_room_id, ok := hostelStudentData["hostel_room_id"]; ok {
		hs.HostelRoomId = uint(hostel_room_id.(int64))
	}

	if hostel_id, ok := hostelStudentData["hostel_id"]; ok {
		hs.HostelId = uint(hostel_id.(int64))
	}

	if contactNumber, ok := hostelStudentData["content_number"]; ok {
		hs.ContactNumber = contactNumber.(string)
	}

	if feeIncluded, ok := hostelStudentData["fee_included"]; ok {
		hs.FeeIncluded = feeIncluded.(bool)
	}
}

func (hs *HostelStudent) All() ([]HostelStudent, error) {
	var hostelStudents []HostelStudent
	err := db.Driver.Find(&hostelStudents).Error
	return hostelStudents, err
}

func (hs *HostelStudent) Find() error {
	err := db.Driver.First(hs, "ID = ?", hs.ID).Error
	return err
}

func (hs *HostelStudent) Create() error {
	err := db.Driver.Create(hs).Error
	hs.updateCount()
	if err != nil {
		return err
	} else {
		err = hs.AddTransaction()
	}
	return err
}

func (hs *HostelStudent) updateCount() {
	var count int64
	db.Driver.Model(&HostelStudent{}).Where("hostel_id = ?", hs.HostelId).Count(&count)
	db.Driver.Model(&Hostel{}).Where("id = ?", hs.HostelId).Update("hostel_students_count", count)

	db.Driver.Model(&HostelStudent{}).Where("hostel_room_id = ?", hs.HostelRoomId).Count(&count)
	db.Driver.Model(&HostelRoom{}).Where("id = ?", hs.HostelRoomId).Update("hostel_students_count", count)
}

func (hs *HostelStudent) Update() error {
	err := db.Driver.Save(hs).Error
	return err
}

func (hs *HostelStudent) Delete() error {
	err := db.Driver.Delete(hs).Error
	return err
}

func (hs *HostelStudent) AddTransaction() error {
	hostel := &Hostel{ID: hs.HostelId}
	err := hostel.Find()
	if err != nil {
		return err
	}

	transactionCategory, err := hostel.GetTransactionCategory()
	if err != nil {
	 	return err
	}
	amount := 0.0
	if !hs.FeeIncluded {
		amount = hostel.Rate
	}
	transactionData := map[string]interface{}{"name": "New Hostel Adminission", "student_id": float64(hs.StudentId), 
		"hostel_student_id": float64(hs.ID), "transaction_category_id": float64(transactionCategory.ID),
		"is_cleared": true, "transaction_type": "debit", "amount": amount}

	student := &Student{ID: hs.StudentId}
	err = hostel.Find()
	if err != nil {
		return err
	}

	transaction := NewTransaction(transactionData, *student)
	err = transaction.Create()
	
	return err
}
