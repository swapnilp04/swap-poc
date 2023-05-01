package models

import (
	"fmt"
	"errors"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
)

type Student struct {
	ID            						uint    `json:"id"`
	Inil     									string `json:"inil"`
	FirstName     						string `json:"first_name" validate:"nonzero"`
	MiddleName     						string `json:"middle_name" validate:"nonzero"`
	LastName      						string `json:"last_name" validate:"nonzero"`
	Age           						int    `json:"age" validate:"min=15"`
	ParentName								string `json:"parent_name" validate:"nonzero"`
	ParentOccupation					string `json:"parent_occupation" validate:"nonzero"`
	ContactNumber 						string  `json:"contact_number" gorm:"contact_number" validate:"nonzero,min=10,max=12"`
	WhNumber									string  `json:"wh_number" gorm:"wh_number,min=10,max=12"`
	Status 										string `json:"status"`
	Town 											string `json:"town"`
	HasHostel									bool `json:"has_hostel" gorm:"default:false"`
	BatchStandardStudents     []BatchStandardStudent 
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
	if errs := validator.Validate(s); errs != nil {
	return errs
	} else {
	return nil
	}
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

	if middleName, ok := studentData["middle_name"]; ok {
		s.MiddleName = middleName.(string)
	}

	if lastName, ok := studentData["last_name"]; ok {
		s.LastName = lastName.(string)
	}

	if age, ok := studentData["age"]; ok {
		s.Age = int(age.(float64))
	}
	if parentName, ok := studentData["parent_name"]; ok {
		s.ParentName = parentName.(string)
	}
	if parentOccupation, ok := studentData["parent_occupation"]; ok {
		s.ParentOccupation = parentOccupation.(string)
	}
	if contactNumber, ok := studentData["contact_number"]; ok {
		s.ContactNumber = contactNumber.(string)
	}

	if whNumber, ok := studentData["wh_number"]; ok {
		s.WhNumber = whNumber.(string)
	}

	if town, ok := studentData["town"]; ok {
		s.Town = town.(string)
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
		
		batchStandardStudent.BatchId = batchStandard.BatchId
		batchStandardStudent.StandardId = batchStandard.StandardId
		batchStandardStudent.StudentId = s.ID
		batchStandardStudent.BatchStandardId = batchStandard.ID
		batchStandardStudent.Fee = batchStandard.Fee
		
		return batchStandardStudent.Create()
	}
	
} 

func (s *Student) AssignHostel(h *Hostel, hr *HostelRoom) error {
	var hostelStudent = HostelStudent{StudentId: s.ID, HostelId: h.ID, HostelRoomId: hr.ID}
	err := db.Driver.Where("student_id = ?",s.ID).First(&hostelStudent).Error
	
	if err != nil {
		err = hostelStudent.Create()
		if err == nil {
			s.HasHostel = true
			s.Update()
		}
	}
	return err
}

func (s *Student) ChangeHostel(h *Hostel, hr *HostelRoom) error {
	var hostelStudent = HostelStudent{StudentId: s.ID}
	err := db.Driver.Where("student_id = ?",s.ID).First(&hostelStudent).Error
	if err == nil {
		hostelStudent.HostelId = h.ID
		hostelStudent.HostelRoomId = hr.ID
		s.HasHostel = true
		s.Update()
		return hostelStudent.Update()
	}
	return err
}

func (s *Student) GetStudentHostel() (HostelStudent, error) {
	var hostelStudent = HostelStudent{}
	err := db.Driver.Where("student_id = ?", s.ID).Preload("Hostel").Preload("HostelRoom").First(&hostelStudent).Error
	
	return hostelStudent, err
}

func (s *Student) GetTransactions() ([]Transaction, error) {
	transactions := []Transaction{}
	err := db.Driver.Where("StudentId = ?", s.ID).Find(transactions).Error
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

func (s *Student) PayCash(amount float64, parentName string ) error {
	transaction := &Transaction{}
	
	transactionData := map[string]interface{}{"Name": "Pay fee", "StudentId": s.ID, 
		"IsCleared": true, "TransactionType": "cridit", "PaidBy": parentName, "PaymentMode": "Cash",
		"Amount": amount}
	transaction.Assign(transactionData)
	err := transaction.Create()

	return err
}

func (s *Student) PayCheque(amount float64, chequeNo string, bankName string, date string) error {
	transaction := &Transaction{}

	transactionData := map[string]interface{}{"Name": "Pay fee", "StudentId": s.ID, 
		"IsCleared": true, "TransactionType": "cridit", 
		"Amount": amount}
	transaction.Assign(transactionData)
	err := transaction.Create()
	if err == nil {
		cheque := &Cheque{}
		chequeData := map[string]interface{}{"BankName": bankName, "Amount": amount, "TransactionID": transaction.ID, "IsCleared": false, "Date": date}
		cheque.Assign(chequeData)
		err = cheque.Create()
	}
	return err
}

