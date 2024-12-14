package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetTeacherLogs(c echo.Context) error {
	tl := &models.TeacherLog{}
	teacherLogs, err := tl.All()
	if err != nil {
		fmt.Println("s.ALL(GetTeacherLogs)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, teacherLogs)
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
	err := teacher.FindByUser
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	teacherLog := models.NewTeacherLog(teacherLogData)
	teacherLog.TeacherID = teacher.ID
	if err := teacherLog.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	if err := teacherLog.Create(); err != nil {
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

	tl.Assign(teacherLogData)
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "TeacherLog deleted successfully"})
}
