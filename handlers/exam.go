package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetExams(c echo.Context) error {
	// Get all users
	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}

	e := &models.Exam{}
	exams, err := e.All(newPage)
	if err != nil {
		fmt.Println("s.ALL(GetExams)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	count, err := e.AllCount()

	return c.JSON(http.StatusOK, map[string]interface{}{"exams": exams, "total": count})
}

func GetExam(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	e := &models.Exam{ID: uint(newId)}
	err = e.Find()
	if err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, e)
}

func CreateExam(c echo.Context) error {
	examData := make(map[string]interface{})
	if err := c.Bind(&examData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("exams %+v\n", examData)
	exam := models.NewExam(examData)
	err := exam.AssighBatchStandard()
	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	if err := exam.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err = exam.Create()
	if err != nil {
		
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	examChapters := examData["exam_chapters"].([]interface {})
	err = exam.AssignExamChapters(examChapters)
	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "exam created", "exam": exam})
}

func UpdateExam(c echo.Context) error {
	cc := c.(CustomContext)
	session := cc.session
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	examData := make(map[string]interface{})

	if err := c.Bind(&examData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	e := &models.Exam{ID: uint(newId)}
	if err := e.Find(); err != nil {
		fmt.Println("s.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	if session.Role != "Admin"  {
		if e.ExamStatus != "Created" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Can not update Exam"})
		}
	}

	e.Assign(examData)
	err = e.AssighBatchStandard()

	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	if err := e.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}
	if err := e.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	examChapters := examData["exam_chapters"].([]interface {})
	err = e.AssignExamChapters(examChapters)
	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "exam updated", "exam": e})
}

func DeleteExam(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	e := &models.Exam{ID: uint(newId)}
	if err := e.Find(); err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := e.Delete(); err != nil {
		fmt.Println("s.Delete(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "exam deleted successfully"})
}


func ConductExam(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	e := &models.Exam{ID: uint(newId)}
	if err := e.Find(); err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	if e.ExamStatus != "Created" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Can not Conduct Exam"})
	}
	
	if err := e.PlotExamStudents(); err != nil {
		fmt.Println("s.ConductExam()", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "exam Conducted", "exam": e})	
}

func PublishExam(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	e := &models.Exam{ID: uint(newId)}
	if err := e.Find(); err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	if e.ExamStatus != "Conducted" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Can not Conduct Exam"})
	}
	
	if err := e.PublishExam(); err != nil {
		fmt.Println("s.ConductExam()", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "exam Published", "exam": e})	
}

func GetExamStudents(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	e := &models.Exam{ID: uint(newId)}
	if err := e.Find(); err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	examStudents, err := e.GetExamStudents()
	if err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, examStudents)
}

func SaveExamMarks(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	e := &models.Exam{ID: uint(newId)}
	if err := e.Find(); err != nil {
		fmt.Println("e.Find(GetExam)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if e.ExamStatus != "Conducted" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Exam Can not be on Conducted State"})
	}

	//examStudents := make([]interface{}, 0)
	examStudents := make([]map[string]interface{}, 0)
	
	if err := c.Bind(&examStudents); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}	
	
	err = e.SaveExamMarks(examStudents)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func GetExamGroupReport(c echo.Context) error { 
	examIds := c.QueryParam("examString")
	exam := &models.Exam{}
	examStudents, err := exam.GetExamsReportStudents(examIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	return c.JSON(http.StatusOK, examStudents)
}