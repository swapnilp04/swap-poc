package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Hostel struct {
	ID            	uint    `json:"id"`
	Name     				string `json:"name"`
	Rooms      			int `json:"rooms"`
	Rector      		string `json:"rector"`	
	ContactNumber 	int64  `json:"contact_number" gorm:"contact_number"`
	Rate     				float64 	`json:"rate"` 
	HostelRoomsCount int64 `json:"hostel_rooms_count"`
	HostelStudentsCount int64 `json:"hostel_students_count"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateHostel() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&Hostel{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewHostel(hostelData map[string]interface{}) *Hostel {
	hostel := &Hostel{}
	hostel.Assign(hostelData)
	return hostel
}

func (s *Hostel) Validate() error {
	return nil
}

func (s *Hostel) HostelRooms() ([]HostelRoom, error){
	var hostelRooms []HostelRoom
	err := db.Driver.Where("hostel_id = ?", s.ID).Find(&hostelRooms).Error
	return hostelRooms, err
}

func (h *Hostel) Assign(hostelData map[string]interface{}) {
	fmt.Printf("%+v\n", hostelData)
	if name, ok := hostelData["name"]; ok {
		h.Name = name.(string)
	}

	if rooms, ok := hostelData["rooms"]; ok {
		h.Rooms = int(rooms.(float64))
	}

	if rate, ok := hostelData["rate"]; ok {
		h.Rate = rate.(float64)
	}

	if rector, ok := hostelData["rector"]; ok {
		h.Rector = rector.(string)
	}


	if contactNumber, ok := hostelData["content_number"]; ok {
		h.ContactNumber = int64(contactNumber.(float64))
	}
}

func (h *Hostel) All() ([]Hostel, error) {
	var hostels []Hostel
	err := db.Driver.Find(&hostels).Error
	return hostels, err
}

func (h *Hostel) Find() error {
	err := db.Driver.First(h, "ID = ?", h.ID).Error
	return err
}

func (h *Hostel) Create() error {
	err := db.Driver.Create(h).Error
	if err != nil {
		return err
	} else {
		err = h.createTransactionCategory()
	}
	return err
}

func (h *Hostel) Update() error {
	err := db.Driver.Save(h).Error
	return err
}

func (h *Hostel) Delete() error {
	err := db.Driver.Delete(h).Error
	return err
}

func (h *Hostel) createTransactionCategory() error {
	var transactionCatetoryData = map[string]interface{}{"name": "Hostel", "hostel_id": h.ID}
	transactionCategory := NewTransactionCategory(transactionCatetoryData)
	err := transactionCategory.Create()
	return err
}

func (h *Hostel) GetTransactionCategory() (*TransactionCategory, error) {
	tc := &TransactionCategory{Name: "Hostel", HostelId: h.ID}
	err := db.Driver.Where("name like ? and hostel_id = ?", "Hostel", h.ID).First(tc).Error
	return tc, err
}