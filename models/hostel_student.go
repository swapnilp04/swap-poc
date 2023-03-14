package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type HostelStudent struct {
	ID            	int    `json:"id"`
	Name     				string `json:"name"`
	HostelId				int `json:"hostel_id"`
	RoomId      		int `json:"room_id"`
	ContactNumber  string `json:"contact_number"`
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

	if room_id, ok := hostelStudentData["room_id"]; ok {
		hs.RoomId = int(room_id.(int64))
	}

	if hostel_id, ok := hostelStudentData["hostel_id"]; ok {
		hs.HostelId = int(hostel_id.(int64))
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
	return err
}

func (hs *HostelStudent) Update() error {
	err := db.Driver.Save(hs).Error
	return err
}

func (hs *HostelStudent) Delete() error {
	err := db.Driver.Delete(hs).Error
	return err
}
