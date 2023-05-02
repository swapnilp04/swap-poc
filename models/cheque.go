package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Cheque struct {
	ID            					uint    `json:"id"`
	BankName     						string `json:"bank_name"`
	IsCleared 							bool `json:"is_cleared"`
	Amount       						float64 `json:"amount"`
	TransactionId 					uint `json:"transaction_id"`
	Date  									time.Time
	CreatedAt 							time.Time
	UpdatedAt 							time.Time
  DeletedAt 							gorm.DeletedAt `gorm:"index"`
}


func migrateCheque() {
	fmt.Println("migrating Cheque..")
	err := db.Driver.AutoMigrate(&Cheque{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewCheque(chequeData map[string]interface{}) *Cheque {
	cheque := &Cheque{}
	cheque.Assign(chequeData)
	return cheque
}

func (c *Cheque) Validate() error {
	return nil
}

func (c *Cheque) Assign(chequeData map[string]interface{}) {
	fmt.Printf("%+v\n", chequeData)
	if bankName, ok := chequeData["bank_name"]; ok {
		c.BankName = bankName.(string)
	}

	if isCleared, ok := chequeData["is_cleared"]; ok {
		c.IsCleared = isCleared.(bool)
	}	

	if amount, ok := chequeData["amount"]; ok {
		c.Amount = amount.(float64)
	}	

	if date, ok := chequeData["date"]; ok {
		c.Date, _ = time.Parse("2006-01-02T15:04:05.999999999Z", date.(string))
	}	
}

func (c *Cheque) All() ([]Cheque, error) {
	var cheques []Cheque
	err := db.Driver.Find(&cheques).Error
	return cheques, err
}

func (c *Cheque) Find() error {
	err := db.Driver.First(c, "ID = ?", c.ID).Error
	return err
}

func (c *Cheque) Create() error {
	err := db.Driver.Create(c).Error
	return err
}

func (c *Cheque) Update() error {
	err := db.Driver.Save(c).Error
	return err
}

func (c *Cheque) Delete() error {
	err := db.Driver.Delete(c).Error
	return err
}
