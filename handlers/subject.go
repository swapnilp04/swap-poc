package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetSubjects(c echo.Context) error {
	// Get all users
	id := c.Param("standard_id")
	standard, err := findStandard(id)
	if err != nil {
		fmt.Println("s.Find(GetStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	subjects, err := standard.GetSubjects()
	
	if err != nil {
		fmt.Println("s.ALL(GetSubjects)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, subjects)
}

func GetSubject(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	s := &models.Subject{ID: uint(newId)}
	err = s.Find()
	if err != nil {
		fmt.Println("s.Find(GetSubject)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func CreateSubject(c echo.Context) error {
	id := c.Param("standard_id")
	standard, err := findStandard(id)
	if err != nil {
		fmt.Println("s.Find(GetStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	subjectData := make(map[string]interface{})
	if err := c.Bind(&subjectData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}
	
	subject := models.NewSubject(subjectData)
	subject.StandardID = standard.ID
	if err := subject.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err = subject.Create()
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "subject created", "subject": subject})
}

func UpdateSubject(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	subjectData := make(map[string]interface{})

	if err := c.Bind(&subjectData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	s := &models.Subject{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetSubject)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	s.Assign(subjectData)
	if err := s.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Subject updated", "subject": s})
}

func DeleteSubject(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	s := &models.Subject{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetSubject)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := s.Delete(); err != nil {
		fmt.Println("s.Delete(GetSubject)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Subject deleted successfully"})
}

func findStandard(id string) (models.Standard, error){
	newId, err := strconv.Atoi(id)
	s := models.Standard{ID: uint(newId)}
	if err != nil {
		return s, err
	}

	err = s.Find()
	return s, err
}
