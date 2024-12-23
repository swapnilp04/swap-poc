package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetTeachers(c echo.Context) error {
	// Get all users

	teacher := &models.Teacher{}
	teachers, err := teacher.All()
	if err != nil {
		fmt.Println("s.ALL(GetTeachers)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	return c.JSON(http.StatusOK, teachers)
}

func GetTeacher(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	teacher := &models.Teacher{ID: uint(newId)}
	err = teacher.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, teacher)
}

func CreateTeacher(c echo.Context) error {
	teacherData := make(map[string]interface{})
	if err := c.Bind(&teacherData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("teachers %+v\n", teacherData)
	teacher := models.NewTeacher(teacherData)
	if err := teacher.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err := teacher.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	err = teacher.CreateUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "teacher created", "teacher": teacher})
}

func UpdateTeacher(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	teacherData := make(map[string]interface{})

	if err := c.Bind(&teacherData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	t := &models.Teacher{ID: uint(newId)}
	if err := t.Find(); err != nil {
		fmt.Println("s.Find(GetTeacher)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	t.Assign(teacherData)
	if err := t.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}
	if err := t.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "teacher updated", "teacher": t})
}

func DeleteTeacher(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	t := &models.Teacher{ID: uint(newId)}
	if err := t.Find(); err != nil {
		fmt.Println("s.Find(GetTeacher)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := t.Delete(); err != nil {
		fmt.Println("s.Delete(GetTeacher)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "teacher deleted successfully"})
}


func GetTeachersLog(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	teacher := &models.Teacher{ID: uint(newId)}
	err = teacher.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}
	searchBatchStandard := c.QueryParam("searchBatchStandard")
	searchSubject := c.QueryParam("searchSubject")
	searchDate := c.QueryParam("searchDate")
	
	teacherLogs, err := teacher.GetTeachersLogs(newPage, searchBatchStandard, searchSubject, searchDate)
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	count, err := teacher.AllTeachersLogsCount(searchBatchStandard, searchSubject, searchDate)

	return c.JSON(http.StatusOK, map[string]interface{}{"teacherLogs": teacherLogs, "total": count})
}

