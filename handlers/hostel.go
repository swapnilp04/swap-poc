package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetHostels(c echo.Context) error {
	// Get all users
	b := &models.Hostel{}
	hostels, err := b.All()
	if err != nil {
		fmt.Println("s.ALL(GetHostels)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, hostels)
}

func GetHostel(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	b := &models.Hostel{ID: uint(newId)}
	err = b.Find()
	if err != nil {
		fmt.Println("s.Find(GetHostel)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, b)
}

func CreateHostel(c echo.Context) error {
	hostelData := make(map[string]interface{})
	if err := c.Bind(&hostelData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("hostels %+v\n", hostelData)
	hostel := models.NewHostel(hostelData)
	if err := hostel.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	err := hostel.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "hostel created", "hostel": hostel})
}

func UpdateHostel(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	hostelData := make(map[string]interface{})

	if err := c.Bind(&hostelData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	s := &models.Hostel{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetHostel)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	s.Assign(hostelData)
	if err := s.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Hostel updated", "hostel": s})
}

func DeleteHostel(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	b := &models.Hostel{ID: uint(newId)}
	if err := b.Find(); err != nil {
		fmt.Println("b.Find(GetHostel)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := b.Delete(); err != nil {
		fmt.Println("b.Delete(GetHostel)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Hostel deleted successfully"})
}
