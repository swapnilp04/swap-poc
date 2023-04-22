package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetBatchs(c echo.Context) error {
	// Get all users
	b := &models.Batch{}
	batchs, err := b.All()
	if err != nil {
		fmt.Println("s.ALL(GetBatchs)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchs)
}

func GetBatch(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	b := &models.Batch{ID: uint(newId)}
	err = b.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, b)
}

func CreateBatch(c echo.Context) error {
	batchData := make(map[string]interface{})
	if err := c.Bind(&batchData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("batchs %+v\n", batchData)
	batch := models.NewBatch(batchData)
	if err := batch.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	err := batch.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "batch created", "batch": batch})
}

func UpdateBatch(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	batchData := make(map[string]interface{})

	if err := c.Bind(&batchData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	s := &models.Batch{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	s.Assign(batchData)
	if err := s.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Batch updated", "batch": s})
}

func DeleteBatch(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	b := &models.Batch{ID: uint(newId)}
	if err := b.Find(); err != nil {
		fmt.Println("b.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := b.Delete(); err != nil {
		fmt.Println("b.Delete(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Batch deleted successfully"})
}
