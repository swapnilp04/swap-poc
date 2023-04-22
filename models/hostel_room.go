package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type HostelRoom struct {
	ID            	uint    	`json:"id"`
	Name     				string 	`json:"name"`
	NoOfStudents    int 		`json:"no_of_students"`
	Rate     				int64 	`json:"rate"`
	HostelID        uint `json:"hostel_id"`
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
	return nil
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
		hr.Rate = int64(rate.(float64))
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
	return err
}

func (hr *HostelRoom) Update() error {
	err := db.Driver.Save(hr).Error
	return err
}

func (hr *HostelRoom) Delete() error {
	err := db.Driver.Delete(hr).Error
	return err
}
