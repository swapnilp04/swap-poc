package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type HostelRoom struct {
	ID            	uint    	`json:"id"`
	Name     				string 	`json:"name" validate:"nonzero"`
	NoOfStudents    int 		`json:"no_of_students"`
	Rate     				float64 	`json:"rate" validate:"nonzero"`
	HostelID        uint `json:"hostel_id" validate:"nonzero"`
	HostelStudentsCount int64 `json:"hostel_students_count"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateHostelRoom() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&HostelRoom{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewHostelRoom(hostelRoomData map[string]interface{}) *HostelRoom {
	hostelRoom := &HostelRoom{}
	hostelRoom.Assign(hostelRoomData)
	return hostelRoom
}

func (hr *HostelRoom) Validate() error {
	if errs := validator.Validate(hr); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (hr *HostelRoom) Assign(hostelRoomData map[string]interface{}) {
	fmt.Printf("%+v\n", hostelRoomData)
	if name, ok := hostelRoomData["name"]; ok {
		hr.Name = name.(string)
	}

	if noOfStudents, ok := hostelRoomData["no_of_students"]; ok {
		hr.NoOfStudents = int(noOfStudents.(float64))
	}

	if rate, ok := hostelRoomData["rate"]; ok {
		hr.Rate = rate.(float64)
	}
}

func (hr *HostelRoom) All() ([]HostelRoom, error) {
	var hostelRooms []HostelRoom
	err := db.Driver.Find(&hostelRooms).Error
	return hostelRooms, err
}

func (hr *HostelRoom) Find() error {
	err := db.Driver.First(hr, "ID = ?", hr.ID).Error
	return err
}

func (hr *HostelRoom) Create() error {
	err := db.Driver.Create(hr).Error
	hr.updateCount()
	return err
}

func (hr *HostelRoom) updateCount() {
	var count int64
	db.Driver.Model(&HostelRoom{}).Where("hostel_id = ?", hr.HostelID).Count(&count)
	db.Driver.Model(&Hostel{}).Where("id = ?", hr.HostelID).Update("hostel_rooms_count", count)
}

func (hr *HostelRoom) Update() error {
	err := db.Driver.Save(hr).Error
	return err
}

func (hr *HostelRoom) Delete() error {
	err := db.Driver.Delete(hr).Error
	return err
}

func (hr *HostelRoom) GetHostelRoomStudents() ([]HostelStudent, error) {
	var hostelStudents []HostelStudent
	err := db.Driver.Where("hostel_id = ? and hostel_room_id = ?", hr.HostelID, hr.ID).Preload("Student").
	Preload("Hostel").Find(&hostelStudents).Error
	return hostelStudents, err	
}
