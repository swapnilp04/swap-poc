package models

import (
	"fmt"
	"errors"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Student struct {
	ID            						int    `json:"id"`
	Inil     									string `json:"inil"`
	FirstName     						string `json:"first_name"`
	MiddleName     						string `json:"middle_name"`
	LastName      						string `json:"last_name"`
	Age           						int    `json:"age"`
	ParentName								string `json:"parent_name"`
	ParentOccupation					string `json:"parent_occupation"`
	ContactNumber 						int64  `json:"phone_number" gorm:"phone_number"`
	Status 										string `json:"status"`
	/*BatchStandardStudents  		[]BatchStandardStudent
	Transactions  						[]Transaction
	HostelStudent 						HostelStudent
	HostelStudentTransactions []HostelStudentTransaction
	ExamStudents							[]ExamStudent*/
	CreatedAt 								time.Time
	UpdatedAt 								time.Time
  DeletedAt 								gorm.DeletedAt `gorm:"index"`
}

func migrateStudent() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&Student{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewStudent(studentData map[string]interface{}) *Student {
	student := &Student{}
	student.Assign(studentData)
	return student
}

func (s *Student) Validate() error {
	return nil
}

func (s *Student) AssignClass() error {
	
	if s.Status == "Admission" {
		return nil
	} else {
		return errors.New("Already assigned Class")
	}
}

func (s *Student) RemoveFromClass() error {
	return nil
}

func (s *Student) Assign(studentData map[string]interface{}) {
	fmt.Printf("%+v\n", studentData)
	if firstName, ok := studentData["first_name"]; ok {
		s.FirstName = firstName.(string)
	}

	if lastName, ok := studentData["last_name"]; ok {
		s.LastName = lastName.(string)
	}

	if age, ok := studentData["age"]; ok {
		s.Age = int(age.(float64))
	}

	if contactNumber, ok := studentData["content_number"]; ok {
		s.ContactNumber = int64(contactNumber.(float64))
	}
}

func (s *Student) All() ([]Student, error) {
	var students []Student
	err := db.Driver.Find(&students).Error
	return students, err
}

func (s *Student) Find() error {
	err := db.Driver.First(s, "ID = ?", s.ID).Error
	return err
}

func (s *Student) Create() error {
	err := db.Driver.Create(s).Error
	return err
}

func (s *Student) Update() error {
	err := db.Driver.Save(s).Error
	return err
}

func (s *Student) Delete() error {
	err := db.Driver.Delete(s).Error
	return err
}
