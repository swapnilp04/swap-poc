package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetBatchStandard(c echo.Context) error {
	// Get a single user by ID
	batchId := c.Param("batch_id")
	newBatchId, err := strconv.Atoi(batchId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	b := &models.Batch{ID: uint(newBatchId)}
	err = b.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	batchStandardId := c.Param("id")
	newBatchStandardId, err := strconv.Atoi(batchStandardId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	batchStandard := &models.BatchStandard{ID: uint(newBatchStandardId)}
	err = batchStandard.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchStandard)
}

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
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err = batchStandard.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
		
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "batch Standard created", "batch_standard": batchStandard})	
}

func UpdateBatchStandard(c echo.Context) error {
	batchId := c.Param("batch_id")
	newBatchId, err := strconv.Atoi(batchId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	b := &models.Batch{ID: uint(newBatchId)}
	if err := b.Find(); err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}


	batchStandardData := make(map[string]interface{})
	if err := c.Bind(&batchStandardData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}
	
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	bs := &models.BatchStandard{ID: uint(newId)}
	if err := bs.Find(); err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	bs.Assign(batchStandardData)
	if err := bs.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Batch Standard updated", "batch_standard": bs})
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
	standardIds, err := bs.AllIds(b.ID)
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
