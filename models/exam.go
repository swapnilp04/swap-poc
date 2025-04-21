package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	//"strconv"
	"strings"
)

type Exam struct {
	ID            	uint    `json:"id"`
	Name     				string `json:"name" validate:"nonzero"`
	BatchStandardID	uint `json:"batch_standard_id" validate:"nonzero"`
	StandardID      uint `json:"standard_id" validate:"nonzero"`
	Standard 				Standard `validate:"-"`
	BatchID      		uint `json:"batch_id" validate:"nonzero"`
	Batch 					Batch `validate:"-"`
	ExamType				string `json:"exam_type" validate:"nonzero"`
	ExamMarks				int `json:"exam_marks" validate:"nonzero"`
	ExamTime				int `json:"exam_time" validate:"nonzero"`
	ExamDate				time.Time `json:"exam_date"`
	ExamStatus 			string `json:"exam_status" validate:"nonzero" gorm:"default:'Created'"` // Created, Conducted, Published
	SubjectID 			uint `json:"subject_id" validate:"nonzero"`
	Subject 				Subject `validate:"-"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateExam() {
	fmt.Println("migrating Exam..")
	err := db.Driver.AutoMigrate(&Exam{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewExam(examData map[string]interface{}) *Exam {
	exam := &Exam{}
	exam.Assign(examData)
	exam.ExamStatus = "Created"
	return exam
}

func (e *Exam) Validate() error {
	if errs := validator.Validate(e); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (e *Exam) Assign(examData map[string]interface{}) {
	fmt.Printf("%+v\n", examData)
	if name, ok := examData["name"]; ok {
		e.Name = name.(string)
	}
	if batchStandardID, ok := examData["batch_standard_id"]; ok {
		e.BatchStandardID = uint(batchStandardID.(float64))
	}
	if subjectID, ok := examData["subject_id"]; ok {
		e.SubjectID = uint(subjectID.(float64))
	}
	if examDate, ok := examData["exam_date"]; ok {
		e.ExamDate, _ = time.Parse("2006-01-02T15:04:05.999999999Z", examDate.(string))
	}
	if examType, ok := examData["exam_type"]; ok {
		e.ExamType = examType.(string)
	}
	if examMarks, ok := examData["exam_marks"]; ok {
		e.ExamMarks = int(examMarks.(float64))
	}
	if examTime, ok := examData["exam_time"]; ok {
		e.ExamTime = int(examTime.(float64))
	}
}

func (e *Exam) All(page int) ([]Exam, error) {
	var exams []Exam
	err := db.Driver.Limit(10).Preload("Standard").Preload("Subject").Preload("Batch").Offset((page-1) * 10).Order("id desc").Find(&exams).Error
	return exams, err
}

func (e *Exam) AllCount() (int64, error) {
	var count int64
	query := db.Driver.Model(&Exam{})
	err := query.Count(&count).Error
	return count, err
}

func (e *Exam) Find() error {
	err := db.Driver.Preload("Standard").Preload("Subject").Preload("Batch").First(e, "ID = ?", e.ID).Error
	return err
}

func (e *Exam) Create() error {
	err := db.Driver.Create(e).Error
	return err
}

func (e *Exam) Update() error {
	//err := db.Driver.Updates(e).Error
	err := db.Driver.Omit("Subject, Standard").Session(&gorm.Session{FullSaveAssociations: false}).Updates(&e).Error
	return err
}

func (e *Exam) Delete() error {
	err := db.Driver.Delete(e).Error
	return err
}

func(e *Exam) GetBatchStandard() (BatchStandard, error) {
	batchStandard := BatchStandard{}
	
	err := db.Driver.First(&batchStandard, "ID = ?", e.BatchStandardID).Error
	return batchStandard, err
}

func (e * Exam) ChangeStatus(status string) error {
	err := db.Driver.Model(e).Updates(Exam{ExamStatus: status}).Error
	return err
}

func (e *Exam) AssighBatchStandard() error {
	batchStandard, err := e.GetBatchStandard()
	if err == nil {
		e.BatchID = batchStandard.BatchId
		e.StandardID = batchStandard.StandardId
	}
	return err
}

func (e *Exam) PlotExamStudents() error {
	batchStandard, err := e.GetBatchStandard()
	if(err != nil) {
		return err
	}

	batchStandardStudents, err := batchStandard.GetStudents()
	for _, batchStandardStudent := range batchStandardStudents {
			examStudent := &ExamStudent{StudentID: batchStandardStudent.StudentId, ExamID: e.ID}
			examStudent.Create()
		}
	//err := db.Driver.Find(&examStudents).Error
	err = e.ChangeStatus("Conducted")
	return  err
}


func (e *Exam) PublishExam() error {
	err := db.Driver.Model(&e).Updates(Exam{ExamStatus: "Published"}).Error
	return err
}

func (e *Exam) GetExamStudents() ([]ExamStudent, error) {
	var examStudents []ExamStudent
	err := db.Driver.Preload("Student").Where("exam_id = ?", e.ID).Find(&examStudents).Error
	return examStudents, err
}

func (e *Exam) SaveExamMarks(examStudents []map[string]interface{}) error {
	for _, examStudentObj := range examStudents {
		newId := examStudentObj["id"].(float64)		
		es := ExamStudent{ID: uint(newId)}	
		es.Assign(examStudentObj)
		err := es.UpdateMarks()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exam) GetExamsReportStudents(examIds string) ([]ExamStudent, error) {
	var examStudents []ExamStudent
	examsArr := strings.Split(examIds, ",")
	err := db.Driver.Preload("Student").Where("exam_id in (?)", examsArr).Find(&examStudents).Error
	return examStudents, err
}
