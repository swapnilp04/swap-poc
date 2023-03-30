package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type HostelStudentTransaction struct {
	ID            	int    `json:"id"`
	Name     				string `json:"name"`
	HostelId				int `json:"hostel_id"`
	RoomId      		int `json:"room_id"`
	ContactNumber  	string `json:"contact_number"`
	StudentId  			uint `json:"student_id"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateHostelStudentTransaction() {
	fmt.Println("migrating hostel student Transaction..")
	err := db.Driver.AutoMigrate(&HostelStudentTransaction{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewHostelStudentTransaction(hostelStudentTransactionData map[string]interface{}) *HostelStudentTransaction {
	hostelStudentTransaction := &HostelStudentTransaction{}
	hostelStudentTransaction.Assign(hostelStudentTransactionData)
	return hostelStudentTransaction
}

func (hs *HostelStudentTransaction) Validate() error {
	return nil
}

func (hs *HostelStudentTransaction) Assign(hostelStudentTransactionData map[string]interface{}) {
	fmt.Printf("%+v\n", hostelStudentTransactionData)
	if name, ok := hostelStudentTransactionData["name"]; ok {
		hs.Name = name.(string)
	}

	if room_id, ok := hostelStudentTransactionData["room_id"]; ok {
		hs.RoomId = int(room_id.(int64))
	}

	if hostel_id, ok := hostelStudentTransactionData["hostel_id"]; ok {
		hs.HostelId = int(hostel_id.(int64))
	}

	if contactNumber, ok := hostelStudentTransactionData["content_number"]; ok {
		hs.ContactNumber = contactNumber.(string)
	}
}

func (hs *HostelStudentTransaction) All() ([]HostelStudentTransaction, error) {
	var hostelStudentTransactions []HostelStudentTransaction
	err := db.Driver.Find(&hostelStudentTransactions).Error
	return hostelStudentTransactions, err
}

func (hs *HostelStudentTransaction) Find() error {
	err := db.Driver.First(hs, "ID = ?", hs.ID).Error
	return err
}

func (hs *HostelStudentTransaction) Create() error {
	err := db.Driver.Create(hs).Error
	return err
}

func (hs *HostelStudentTransaction) Update() error {
	err := db.Driver.Save(hs).Error
	return err
}

func (hs *HostelStudentTransaction) Delete() error {
	err := db.Driver.Delete(hs).Error
	return err
}
