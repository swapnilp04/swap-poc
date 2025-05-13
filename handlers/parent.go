package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"
	"github.com/labstack/echo/v4"
)

func GetParents(c echo.Context) error {
	// Get all users
	p := &models.Parent{}

	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}

	search := c.QueryParam("search")

	parents, err := p.AllWithPagination(int(newPage), search)
	if err != nil {
		fmt.Println("s.ALL(GetParents)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	count, err := p.Count(search)
	if err != nil {
		fmt.Println("s.ALL(GetParents)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"parents": parents, "total": count})
}

func GetParent(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	p := &models.Parent{ID: uint(newId)}
	err = p.Find()
	if err != nil {
		fmt.Println("s.Find(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, p)
}

func CreateParent(c echo.Context) error {
	parentData := make(map[string]interface{})
	if err := c.Bind(&parentData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("parents %+v\n", parentData)
	parent := models.NewParent(parentData)
	if err := parent.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}
	err := parent.Create()
	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "parent created", "parent": parent})
}

func UpdateParent(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	parentData := make(map[string]interface{})

	if err := c.Bind(&parentData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	s := &models.Parent{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	s.Assign(parentData)
	if err := s.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}
	if err := s.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "parent updated", "parent": s})
}

func DeleteParent(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	s := &models.Parent{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := s.Delete(); err != nil {
		fmt.Println("s.Delete(GetParent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "parent deleted successfully"})
}