package models

import (
	"fmt"
	"errors"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Student struct {
	ID            						uint    `json:"id"`
	Inil     									string `json:"inil"`
	FirstName     						string `json:"first_name"`
	MiddleName     						string `json:"middle_name"`
	LastName      						string `json:"last_name"`
	Age           						int    `json:"age"`
	ParentName								string `json:"parent_name"`
	ParentOccupation					string `json:"parent_occupation"`
	ContactNumber 						int64  `json:"phone_number" gorm:"phone_number"`
	Status 										string `json:"status"`
	BatchStandardStudents  		[]BatchStandardStudent
	Transactions  						[]Transaction
	HostelStudent 						HostelStudent
	ExamStudents							[]ExamStudent
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
	db.Commit()
	return err
}

func (s *Student) Update() error {
	err := db.Driver.Save(s).Error
	db.Commit()
	return err
}

func (s *Student) Delete() error {
	err := db.Driver.Delete(s).Error
	db.Commit()
	return err
}

func (s *Student) AdmissionStatus() bool {
	return s.Status == "Admission"
}

func (s *Student) ConfirmedStatus() bool {
	return s.Status == "Confirmed"
}

func (s *Student) GetBatchStandardStudents() []BatchStandardStudent {
	return s.BatchStandardStudents
}

func (s *Student) RemoveBatchStandard(batchStandard *BatchStandard) error {
	balance := s.GetBalance()
	if balance > 0.0 {
		return errors.New("Please Clear Balance first")
	}

	batchStandardStudent := &BatchStandardStudent{StudentId: s.ID, BatchStandardId: batchStandard.ID}
	err := db.Driver.Find(batchStandardStudent).Error
	if err != nil {
		return err
	}
	return batchStandardStudent.Delete()
}

func (s *Student) AssignBatchStandard(batchStandard *BatchStandard) error {
	//check student already assign to batch standard
	//if assign remove from current batch standard before this check transaction 
	// after that assign new batch standard
	batchStandardStudents := s.GetBatchStandardStudents()
	if len(batchStandardStudents) > 0 {
		return errors.New("Already assigned to Class")
	} else {
		batchStandardStudent := &BatchStandardStudent{}
		batchStandardStudentData := map[string]interface{}{"BatchId": batchStandard.ID, "StandardId": batchStandard.StandardId, "StudentId": s.ID, 
			"Fee": batchStandard.Fee}
		batchStandardStudent.Assign(batchStandardStudentData)
		return batchStandardStudent.Create()
	}
	
} 

func (s *Student) AssignHostel(h *Hostel, hr *HostelRoom) error {
	var hostelStudent = HostelStudent{StudentId: s.ID, HostelId: h.ID, RoomId: hr.ID}
	err := db.Driver.Find(&hostelStudent).Error
	
	if err != nil {
		err = hostelStudent.Create()
	}
	return err
}

func (s *Student) GetTransactions() ([]Transaction, error) {
	transactions := []Transaction{}
	err := db.Driver.Where("StudentId = ? AND IsCleared = ?", 
		s.ID, true).Find(transactions).Error
	return transactions, err
}

func (s *Student) TotalDebits() float64 {
	transactions, err := s.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "debit" {
				total = total + transaction.Amount
			}
		}
	}
	return total
}

func (s *Student) TotalCridits() float64 {
	transactions, err := s.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "credit" {
				total = total + transaction.Amount
			}
		}
	}
	return total	
}

func (s *Student) GetBalance() float64 {
	transactions, err := s.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "debit" {
				total = total + transaction.Amount
			} else {
				total = total - transaction.Amount
			}
		}
	}
	return total		
}