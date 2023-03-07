package models

import (
	"fmt"
	"swapnil-ex/models/db"
)

type HostelRoom struct {
	ID            	int    `json:"id"`
	Name     		string `json:"name"`
	NoOfStudents      		int `json:"no_of_students"`
}

func migrateHostelRoom() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&Hostel{})
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

func (h *HostelRoom) Assign(hostelRoomData map[string]interface{}) {
	fmt.Printf("%+v\n", hostelRoomData)
	if name, ok := hostelRoomData["name"]; ok {
		h.Name = name.(string)
	}

	if noOfStudents, ok := hostelRoomData["no_of_students"]; ok {
		h.NoOfStudents = int(noOfStudents.(int64))
	}
}

func (h *HostelRoom) All() ([]HostelRoom, error) {
	var hostelRooms []HostelRoom
	err := db.Driver.Find(&hostelRooms).Error
	return hostelRooms, err
}

func (h *HostelRoom) Find() error {
	err := db.Driver.First(h, "ID = ?", h.ID).Error
	return err
}

func (h *HostelRoom) Create() error {
	err := db.Driver.Create(h).Error
	return err
}

func (h *HostelRoom) Update() error {
	err := db.Driver.Save(h).Error
	return err
}

func (h *HostelRoom) Delete() error {
	err := db.Driver.Delete(h).Error
	return err
}
