package models

import (
	"fmt"
	"errors"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	"strings"
	//"github.com/showa-93/go-mask"
)

type Student struct {
	ID            						uint    `json:"id"`
	Inil     									string 	`json:"inil"`
	FirstName     						string 	`json:"first_name" validate:"nonzero"`
	MiddleName     						string 	`json:"middle_name" validate:"nonzero"`
	LastName      						string 	`json:"last_name" validate:"nonzero"`
	RollNumber								string 	`json:"roll_number"`
	BirthDate  								time.Time `json:"birth_date"`
	AdharCard									string 	`json:"adhar_card" gorm:"adhar_card" validate:"nonzero,min=12,max=12"`
	ParentName								string 	`json:"parent_name" validate:"nonzero"`
	ParentOccupation					string 	`json:"parent_occupation" validate:"nonzero"`
	ContactNumber 						string  `json:"contact_number" gorm:"contact_number" validate:"nonzero,min=10,max=12" mask:"filled"`
	WhNumber									string  `json:"wh_number" validate:"nonzero,min=10,max=12" mask:"filled"`
	Status 										string 	`json:"status"`
	Town 											string 	`json:"town" validate:"nonzero"`
	HasHostel									bool 		`json:"has_hostel" gorm:"default:false"`
	Balance 									float64 `json:"balance" gorm:"default:0.0" mask:"random0.0"`
	StudentAccountBalance 		float64 `json:"student_account_balance" gorm:"default:0.0"`
	HostelRoomId    					uint 		`json:"hostel_room_id" validate:"-"`
	StandardId      					uint 		`json:"standard_id"`
	Standard 									Standard `validate:"-"`
	LastPaymentOn							*time.Time `json:"last_payment_on"`
	HasAbsconded							bool		`json:"has_absconded" gorm:"default:false"`
	AbscondedAt								*time.Time
	HasLeft 									bool 		`json:"has_left" gorm:"default:false"`
	LeftAt										*time.Time `json:"left_at"`
	BatchStandardStudents     []BatchStandardStudent `validate:"-"`
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

func migrationForSaveStandardID() error {
	var students []Student
 	err := db.Driver.Preload("BatchStandardStudents").Find(&students).Error
	if err != nil {
		return err
	}
	for _, student := range students {
		batchStandardStudents := student.BatchStandardStudents
		if len(batchStandardStudents) > 0 {
			batchStandardStudent := BatchStandardStudent{}
			db.Driver.Where("student_id = ?", student.ID).Last(&batchStandardStudent)
			db.Driver.Model(&student).Omit("BatchStandardStudents").Update("standard_id", batchStandardStudent.StandardId)
		}
	}
	return nil
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
	if firstName, ok := studentData["first_name"]; ok {
		s.FirstName = firstName.(string)
	}

	if middleName, ok := studentData["middle_name"]; ok {
		s.MiddleName = middleName.(string)
	}

	if lastName, ok := studentData["last_name"]; ok {
		s.LastName = lastName.(string)
	}

	if birthDate, ok := studentData["birth_date"]; ok {
		s.BirthDate, _ = time.Parse("2006-01-02T15:04:05.999999999Z", birthDate.(string))
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

	if adharCard, ok := studentData["adhar_card"]; ok {
		s.AdharCard = adharCard.(string)
	}

	if town, ok := studentData["town"]; ok {
		s.Town = town.(string)
	}
}

func (s *Student) All(page int, search string) ([]Student, error) {
	var students []Student
	query := db.Driver.Limit(10).Preload("Standard").Offset((page - 1) * 10).Order("id desc")
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("first_name like ? or middle_name like ? or last_name like ? OR contact_number like ?", search, search, search, search)
	}
	err := query.Find(&students).Error
	return students, err
}

func (s *Student) SearchStudents(search string) ([]Student, error) {
	var students []Student
	query := db.Driver.Preload("Standard").Order("id desc").Where("has_left = ?", false)
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("first_name like ? or middle_name like ? or last_name like ? OR contact_number like ?", search, search, search, search)
	}
	err := query.Find(&students).Error
	return students, err
}

func (s *Student) Count(search string) (int64, error) {
	var count int64
	query := db.Driver.Model(&Student{})
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("first_name like ? or middle_name like ? or last_name like ? OR contact_number like ?", search, search, search, search)
	}
	err := query.Count(&count).Error
	return count, err
}

func (s *Student) AllForReport() ([]Student, error) {
	var students []Student
	err := db.Driver.Select("id,first_name,middle_name, last_name ").Where("has_left = ?", false).Order("id desc").Find(&students).Error
	return students, err
}


func (s *Student) SearchIds(search string) (error, []uint){
	var ids []uint
	query := db.Driver.Model(&Student{})
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("first_name like ? or middle_name like ? or last_name like ? OR contact_number like ? ", search, search, search, search)
	}
	err := query.Pluck("id", &ids).Error
	return err, ids
}

func (s *Student) Find() error {
	err := db.Driver.Preload("Standard").Preload("BatchStandardStudents").Omit("BatchStandardStudents.Student").Omit("BatchStandardStudents.Batch").Omit("BatchStandardStudents.BatchStandard").Omit("BatchStandardStudents.Standard").First(s, "ID = ?", s.ID).Error
	return err
}

func (s *Student) Create() error {
	err := db.Driver.Create(s).Error
	return err
}

func (s *Student) Update() error {
	err := db.Driver.Omit("BatchStandardStudent").Omit("Standard").Save(s).Error
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

func (s *Student) GetStudentComments(page int) ([]Comment, error) {
	comment := Comment{}
	comments, err := comment.AllByStudent(s.ID, page)
	return comments, err
}

func (s *Student) GetStudentCommentsCount() (int64, error) {
	comment := Comment{}
	var count int64
	count, err := comment.AllByStudentCount(s.ID)
	return count, err
}

func (s *Student) GetStudentLogAttendances(page int) ([]LogAttendance, error) {
	logAttendance := LogAttendance{}
	logAttendances, err := logAttendance.AllByStudent(s.ID, page)
	return logAttendances, err
}

func (s *Student) GetStudentLogAttendanceCount() (int64, error) {
	logAttendance := LogAttendance{}
	var count int64
	count, err := logAttendance.AllByStudentCount(s.ID)
	return count, err
}

func (s *Student) GetBatchStandardStudents() ([]BatchStandardStudent, error) {
	var batchStandardStudents []BatchStandardStudent
	err := db.Driver.Where("student_id = ?", s.ID).Preload("Batch").Preload("Standard").Find(&batchStandardStudents).Error
	return batchStandardStudents, err
}

func (s *Student) GetBatchStandardStudent(batchStandardID uint) ([]BatchStandardStudent, error) {
	var batchStandardStudents []BatchStandardStudent
	err := db.Driver.Where("student_id = ? && batch_standard_id = ?", s.ID, batchStandardID).Preload("Batch").Preload("Standard").
	Find(&batchStandardStudents).Error
	return batchStandardStudents, err
}

func (s *Student) RemoveBatchStandard(batchStandard *BatchStandard) error {
	totalDebits, totalCredits, totalDiscounts := s.GetBalance()
	balance := totalCredits + totalDiscounts - totalDebits
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
	batchStandardStudents, _ := s.GetBatchStandardStudent(batchStandard.ID)
	if len(batchStandardStudents) > 0 {
		return errors.New("Already assigned to Class")
	} else {
		batchStandardStudent := &BatchStandardStudent{}
		
		batchStandardStudent.BatchId = batchStandard.BatchId
		batchStandardStudent.StandardId = batchStandard.StandardId
		batchStandardStudent.StudentId = s.ID
		batchStandardStudent.BatchStandardId = batchStandard.ID
		batchStandardStudent.Fee = batchStandard.Fee
		
		s.StandardId = batchStandard.StandardId
		err := s.Update()
		if err != nil {
			return err
		}
		return batchStandardStudent.Create()
	}
	
} 

func (s *Student) AssignHostel(h *Hostel, hr *HostelRoom, fee_included bool, fee_iteration string) error {
	var hostelStudent = HostelStudent{StudentId: s.ID, HostelId: h.ID, HostelRoomId: hr.ID}
	err := db.Driver.Where("student_id = ?",s.ID).First(&hostelStudent).Error
	hostelStudent.FeeIncluded = fee_included
	hostelStudent.FeeIteration = fee_iteration
	if err != nil {
		err = hostelStudent.Create()
		if err == nil {
			s.HasHostel = true
			s.HostelRoomId = hr.ID
			s.Update()
			s.SaveBalance()
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
		s.HostelRoomId = hr.ID
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

func (s *Student) GetStudentHostelRoommate() ([]HostelStudent, error) {
	hostelStudents := []HostelStudent{}
	err := db.Driver.Where("hostel_room_id = ?", s.HostelRoomId).Preload("Hostel").Preload("HostelRoom").Preload("Student").Find(&hostelStudents).Error
	return hostelStudents, err
}

func (s *Student) GetTransactions() ([]Transaction, error) {
	transactions := []Transaction{}
	err := db.Driver.Where("student_id = ?", s.ID).Find(&transactions).Error
	return transactions, err
}

func (s *Student) GetTransaction(transactionID uint) (Transaction, error) {
	transaction := Transaction{}
	err := db.Driver.Preload("Student.Standard").Where("student_id = ? and id = ?", s.ID, transactionID).Find(&transaction).Error
	return transaction, err
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

func (s *Student) SaveBalance() error{
	debits, credits, discounts := s.GetBalance()
	s.Balance =  credits + discounts - debits

	return s.Update()
}

func(s *Student) UpdateLastPaymentOn() error {
	currentTime := time.Now()
	s.LastPaymentOn = &currentTime
	err := s.Update()
	return err
}

func (s *Student) GetBalance() (float64, float64, float64) {
	transactions, err := s.GetTransactions()
	var totalDebits = 0.0
	var totalCredits = 0.0
	var totalDiscount = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "debit" {
				totalDebits = totalDebits + transaction.Amount
			} else if(transaction.TransactionType == "cridit" && transaction.Name == "Discount") {

				totalDiscount = totalDiscount + transaction.Amount
			} else {
				totalCredits = totalCredits + transaction.Amount
			}

		}
	}
	return totalDebits, totalCredits, totalDiscount
}

func (s *Student) GetStudentAccounts() ([]StudentAccount, error) {
	studentAccounts := []StudentAccount{}
	err := db.Driver.Where("student_id = ?", s.ID).Find(&studentAccounts).Error
	return studentAccounts, err
}

func (s *Student) GetStudentAccountBalance() (float64, float64) {
	studentAccounts, err := s.GetStudentAccounts()
	var totalDebits = 0.0
	var totalCredits = 0.0
	if err == nil {
		for _, studentAccount := range studentAccounts {
			if studentAccount.TransactionType == "debit" {
				totalDebits = totalDebits + studentAccount.Amount
			} else {
				totalCredits = totalCredits + studentAccount.Amount
			}
		}
	}
	return totalDebits, totalCredits
}

func (s *Student) SaveStudentAccountBalance() error{
	debits, credits := s.GetStudentAccountBalance()
	s.StudentAccountBalance =  credits - debits
	return s.Update()
}

func (s *Student) GetUpcommingBirthdays() ([]Student, error) {
	var students []Student
	query := db.Driver.Preload("Standard")
	query = query.Where("DAY(birth_date) >= (DAY(NOW()) - 1) and MONTH(birth_date) = MONTH(NOW())").Order("DAY(birth_date) asc")
	err := query.Find(&students).Error
	return students, err	
}

func (s *Student) GetStudentLogAttendancesCount(page int) ([]LogAttendance, error) {
	var logAttendances []LogAttendance
	err := db.Driver.Limit(10).Preload("Subject").Offset((page - 1) * 10).Where("student_id = ?", s.ID).Order("id desc").Find(&logAttendances).Error
	return logAttendances, err
}

func (s *Student) GetStudentAllExams() ([]ExamStudent, error) {
	var examStudents []ExamStudent
	err := db.Driver.Preload("Exam.Subject").Preload("Exam.Batch").Where("student_id = ?", s.ID).Omit("Student").Order("id desc").Find(&examStudents).Error
	return examStudents, err
}

func (s *Student) GetStudentExams(page int) ([]ExamStudent, error) {
	var examStudents []ExamStudent
	err := db.Driver.Limit(10).Preload("Exam.Subject").Preload("Exam.Standard").Preload("Exam.Batch").Offset((page - 1) * 10).Where("student_id = ?", s.ID).Omit("Student").Order("id desc").Find(&examStudents).Error
	return examStudents, err
}

func (s *Student) GetExamsGraphData(subjectID uint) ([]map[string]interface{}, error) {
	var graphData []map[string]interface{}
	rows, err := db.Driver.Model(&ExamStudent{}).Select("exam_students.percentage, exams.exam_date").
	Joins("left join exams on exam_students.exam_id = exams.id").
	Where("exams.subject_id = ? and exam_students.student_id = ? and exams.exam_status = ?", subjectID, s.ID, "Published").Order("exams.exam_date asc").Rows()
	if err != nil {
		return graphData, err
	}
	defer rows.Close()

	for rows.Next() {
		result := map[string]interface{}{}
		db.Driver.ScanRows(rows, &result)
		date :=  result["exam_date"].(time.Time)
		result["exam_date"] = date.Format("02-Jan-2006") 
		graphData = append(graphData, result)
	}
	
	return graphData, nil	
}


func (s *Student) GetStudentExamsCount() (int64, error) {
	examStudent := ExamStudent{}
	var count int64
	count, err := examStudent.AllByStudentCount(s.ID)
	return count, err
}


func (s *Student) LeftAcademy() error {
	// Remove from classs
	batchStandardStudents, err := s.GetBatchStandardStudents()
	if err != nil {
		return err
	}
	for _, batchStandardStudent := range batchStandardStudents {
		err := batchStandardStudent.Delete()
		if err != nil {
			return err
		}
	}
	// Remove From Hostel
	if s.HasHostel {
		hostelStudent, err := s.GetStudentHostel()
		if err != nil {
	 		return err
	 	}
	 	err = hostelStudent.Delete()
	 	if err != nil {
	 		return err
	 	}
 	}

	// Update Balance 
	s.SaveBalance()
	currentTime := time.Now()
	s.LeftAt = &currentTime
	s.HasLeft = true
	s.HostelRoomId = 0
	s.HasHostel = false
	s.StandardId = 0
	err = s.Update()		
	return err
}

func (s *Student) ReJoinAcademy() error {
	// Update Balance 
	s.LeftAt = nil
	s.HasLeft = false
	err := s.Update()		
	return err
}

func (s *Student) GetMonthlyStudentLogReport(month int, year int)([]LogAttendance, error) {
	var logAttendances []LogAttendance
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0,1,0)
	err := db.Driver.Preload("TeacherLog.BatchStandard.Standard").Preload("TeacherLog.BatchStandard.Batch").
				Preload("TeacherLog.Subject").Preload("TeacherLog.Teacher").Preload("TeacherLog.LogCategory").
				Preload("TeacherLog.Chapter").Joins("left join teacher_logs on teacher_logs.id = log_attendances.teacher_log_id").
				Where("student_id = ?", s.ID).Where("teacher_logs.log_date >= ? and teacher_logs.log_date < ?", startDate, endDate).
				Order("teacher_logs.log_date desc").Find(&logAttendances).Error

	return logAttendances, err
}

func (s *Student) GetMonthlyStudentExamReport(month int, year int)([]ExamStudent, error) {
	var examStudents []ExamStudent
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0,1,0)
	err := db.Driver.Preload("Exam.Standard").Preload("Exam.Subject").Preload("Exam.Batch").Preload("Exam.Teacher").
				Preload("Exam.ExamChapters.Chapter").Joins("left join exams on exam_students.exam_id = exams.id").
				Where("exam_students.student_id = ?", s.ID).Where("exams.exam_date >= ? and exams.exam_date < ?", startDate, endDate).
				Order("exams.exam_date desc").Find(&examStudents).Error
	return examStudents, err
}

// func (s *Student) GetMonthlyStudentLogDurations(month int, year int)([]TeacherDuration, error) {
// 	var teacherDurations []TeacherDuration
// 	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
// 	endDate := startDate.AddDate(0,1,0)
// 	rows, err := db.Driver.Model(&LogAttendance{}).
// 		Select("sum(teacher_logs.duration) as duration, DATE_FORMAT(teacher_logs.log_date, '%Y-%m-%d') as log_date, Count(teacher_logs.id) as count").
// 		Joins("left join teacher_logs on teacher_logs.id = log_attendances.teacher_log_id").
// 		Where("log_attendances.student_id = ?", s.ID).Where("teacher_logs.log_date >= ? and teacher_logs.log_date < ?", startDate, endDate).
// 		Group("DATE_FORMAT(teacher_logs.log_date, '%Y-%m-%d')").Rows()
// 		defer rows.Close()
// 		for rows.Next() {
// 			var teacherDuration TeacherDuration
// 		  db.Driver.ScanRows(rows, &teacherDuration)
// 		  teacherDurations = append(teacherDurations, teacherDuration)
// 		}
// 		return teacherDurations, err
// }
