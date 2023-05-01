package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type HostelStudent struct {
	ID            	uint    `json:"id"`
	Name     				string `json:"name"`
	HostelID				uint `json:"hostel_id"`
	HostelRoomID    uint `json:"hostel_room_id"`
	ContactNumber  	string `json:"contact_number"`
	StudentId				uint `json:"student_id"`
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
	return nil
}

func (hs *HostelStudent) Assign(hostelStudentData map[string]interface{}) {
	fmt.Printf("%+v\n", hostelStudentData)
	if name, ok := hostelStudentData["name"]; ok {
		hs.Name = name.(string)
	}

	if hostel_room_id, ok := hostelStudentData["hostel_room_id"]; ok {
		hs.HostelRoomID = uint(hostel_room_id.(int64))
	}

	if hostel_id, ok := hostelStudentData["hostel_id"]; ok {
		hs.HostelID = uint(hostel_id.(int64))
	}

	if contactNumber, ok := hostelStudentData["content_number"]; ok {
		hs.ContactNumber = contactNumber.(string)
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
	return err
}

func (hs *HostelStudent) updateCount() {
	var count int64
	db.Driver.Model(&HostelStudent{}).Where("hostel_id = ?", hs.HostelID).Count(&count)
	db.Driver.Model(&Hostel{}).Where("id = ?", hs.HostelID).Update("hostel_students_count", count)

	db.Driver.Model(&HostelStudent{}).Where("hostel_room_id = ?", hs.HostelRoomID).Count(&count)
	db.Driver.Model(&HostelRoom{}).Where("id = ?", hs.HostelRoomID).Update("hostel_students_count", count)
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
	transaction := &Transaction{}
	 transactionCategory, err := hs.Hostel.GetTransactionCategory()
	 if err == nil {
	 	return err
	 }
	 transactionData := map[string]interface{}{"Name": "New Adminission", "StudentId": hs.StudentId, "HostelStudentID": hs.ID, 
	 	"TransactionCategoryId": transactionCategory.ID, "IsCleared": true, "TransactionType": "debit", 
	 	"Amount": hs.Hostel.Rate}
	 transaction.Assign(transactionData)
	 err = transaction.Create()
	
	return nil
}
