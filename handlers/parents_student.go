package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"
	"github.com/labstack/echo/v4"
	"github.com/showa-93/go-mask"
)

func GetParentsStudents(c echo.Context) error {
	// Get all users
	parentID := c.Param("parent_id")
	newParentID, err := strconv.Atoi(parentID)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	parent := &models.Parent{ID: uint(newParentID)}
	err = parent.Find()
	if err != nil {
		fmt.Println("s.Find(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	parentStudents, err := parent.GetParentStudents()
	if err != nil {
		fmt.Println("s.ALL(GetParents)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	cc := c.(CustomContext)
	if( !(cc.session.Role == "Admin" || cc.session.Role == "Accountant")) {
		masker := mask.NewMasker()
		masker.SetMaskChar("+")
		parentStudents, _ = mask.Mask(parentStudents)
	}

	return c.JSON(http.StatusOK, parentStudents)
}

func GetParentsStudent(c echo.Context) error {
	// Get all users
	parentID := c.Param("parent_id")
	newParentID, err := strconv.Atoi(parentID)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	parent := &models.Parent{ID: uint(newParentID)}
	err = parent.Find()
	if err != nil {
		fmt.Println("s.Find(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	id := c.Param("id")
	newID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	parentsStudent, err := parent.GetParentStudent(uint(newID))
	if err != nil {
		fmt.Println("s.ALL(GetParents)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}	
	return c.JSON(http.StatusOK, parentsStudent)
}

func CreatParentsStudent(c echo.Context) error {
	parentID := c.Param("parent_id")
	newParentID, err := strconv.Atoi(parentID)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	parent := &models.Parent{ID: uint(newParentID)}
	err = parent.Find()
	if err != nil {
		fmt.Println("s.Find(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	student := &models.Student{ID: uint(newId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	parentStudent, err := parent.AssignStudentToParent(student)
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, parentStudent)
}