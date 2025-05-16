package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"
	"time"
	"github.com/labstack/echo/v4"
)

func GetTeacherLogs(c echo.Context) error {
	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}
	searchBatchStandard := c.QueryParam("searchBatchStandard")
	searchSubject := c.QueryParam("searchSubject")
	searchTeacher := c.QueryParam("searchTeacher")
	searchDate := c.QueryParam("searchDate")

	tl := &models.TeacherLog{}
	teacherLogs, err := tl.All(newPage, searchBatchStandard, searchSubject, searchTeacher, searchDate)
	if err != nil {
		fmt.Println("s.ALL(GetTeacherLogs)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	count, err := tl.AllCount(searchBatchStandard, searchSubject, searchTeacher, searchDate)
	return c.JSON(http.StatusOK, map[string]interface{}{"teacherLogs": teacherLogs, "total": count})
}

func GetTeacherLog(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	tl := &models.TeacherLog{ID: uint(newId)}
	err = tl.Find()
	if err != nil {
		fmt.Println("s.Find(GetTeacherLog)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, tl)
}

func CreateTeacherLog(c echo.Context) error {
	teacherLogData := make(map[string]interface{})
	if err := c.Bind(&teacherLogData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	cc := c.(CustomContext)
	if cc.session == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	session := cc.session
	teacher := models.Teacher{UserID: uint(session.UserID)}
	
	err := teacher.FindByUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	t := time.Now()
	teacherLog := models.NewTeacherLog(teacherLogData)
	teacherLog.TeacherID = teacher.ID
	teacherLog.LogDate = &t
	teacherLog.UserID = uint(session.UserID)
	if err := teacherLog.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	if err := teacherLog.Create(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	

	combinedClasses := teacherLogData["combined_classes"].([]interface {})
	 
	if err :=teacherLog.CreateCombinedClasses(combinedClasses); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "teacherLog created", "teacherLog": teacherLog})
}

func UpdateTeacherLog(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	teacherLogData := make(map[string]interface{})

	if err := c.Bind(&teacherLogData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	tl := &models.TeacherLog{ID: uint(newId)}
	if err := tl.Find(); err != nil {
		fmt.Println("s.Find(GetTeacherLog)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	tl.AssignUpdate(teacherLogData)
	if err := tl.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}
	if err := tl.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "TeacherLog updated", "teacherLog": tl})
}

func DeleteTeacherLog(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	s := &models.TeacherLog{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetTeacherLog)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := s.Delete(); err != nil {
		fmt.Println("s.Delete(GetTeacherLog)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "TeacherLog deleted successfully", "id": newId})
}

func GetLogCategories(c echo.Context) error {
	lcs := &models.LogCategory{}
	logCategories, err := lcs.All()
	if err != nil {
		fmt.Println("s.ALL(GetTeacherLogs)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, logCategories)
}

func GetLogAttendances(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	tl := &models.TeacherLog{ID: uint(newId)}
	if err := tl.Find(); err != nil {
		fmt.Println("s.Find(GetTeacherLog)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	logAttendances, err := tl.GetLogAttendances()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, logAttendances)
}

func ToggleLogAttendance(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	logAttendance := &models.LogAttendance{ID: uint(newId)}
	if err := logAttendance.Find(); err != nil {
		fmt.Println("s.Find(GetTeacherLog)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	err = logAttendance.ToggleAttendance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	 err = logAttendance.TeacherLog.UpdateAbsentsCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Toggle Attendance"})
}
