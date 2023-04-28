package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)


func CreateBatchStandard(c echo.Context) error {
	batchId := c.Param("batch_id")
	newBatchId, err := strconv.Atoi(batchId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	batchStandardData := make(map[string]interface{})
	if err := c.Bind(&batchStandardData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	batch := &models.Batch{ID: uint(newBatchId)}
	if err := batch.Find(); err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}



	fmt.Printf("batchs %+v\n", batchStandardData)
	batchStandard := models.NewBatchStandard(batchStandardData, batch)
	if err := batchStandard.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	err = batchStandard.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "batch Standard created", "batch_standard": batchStandard})
}

func UpdateBatchStandard(c echo.Context) error {
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

func DeleteBatchStandard(c echo.Context) error {
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

func GetBatchStandards(c echo.Context) error {
	batchID := c.Param("batch_id")
	newBatchID, err := strconv.Atoi(batchID)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}	
	b := &models.Batch{ID: uint(newBatchID)}
	err = b.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	bs := &models.BatchStandard{}
	batchStandards, err := bs.All(b.ID)
	if err != nil {
		fmt.Println("s.ALL(GetBatchs)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, batchStandards)
}

func GetBatchUnassignedStandards(c echo.Context) error {
	batchID := c.Param("batch_id")
	newBatchID, err := strconv.Atoi(batchID)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}	
	b := &models.Batch{ID: uint(newBatchID)}
	err = b.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	bs := &models.BatchStandard{}
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	standardIds, err := bs.AllIds(b.ID)
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	if err != nil {
		fmt.Println("s.ALL(GetBatchStandardIds)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	standard := &models.Standard{}
	standards, err := standard.AllExcept(standardIds)
	if err != nil {
		fmt.Println("s.ALL(GetBatchStandards)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}


	return c.JSON(http.StatusOK, standards)
}
