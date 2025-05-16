package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

// func GetBatchStandardStudents(c echo.Context) error {
// 	studentId := c.Param("student_id")
// 	newStudentId, err := strconv.Atoi(studentId)
// 	if err != nil {
// 		fmt.Println("strconv.Atoi failed", err)
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
// 	}
// 	student := &models.Student{ID: uint(newStudentId)}
// 	err = student.Find()
// 	if err != nil {
// 		fmt.Println("s.Find(GetStudent)", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
// 	}
	
// 	bss := &models.BatchStandardStudent{}
// 	batchStandardStudents, err := bss.All(student.ID)
// 	if err != nil {
// 		fmt.Println("s.ALL(GetBatchStandardStudents)", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
// 	}
// 	return c.JSON(http.StatusOK, batchStandardStudents)
//}


func GetBatchStandardStudent(c echo.Context) error {
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	studentBatchStandardId := c.Param("id")
	newStudentBatchStandardId, err := strconv.Atoi(studentBatchStandardId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	batchStandardStudent := &models.BatchStandardStudent{ID: uint(newStudentBatchStandardId)}
	err = batchStandardStudent.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandardStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, batchStandardStudent)
}

func CreateStudentBatchStandard(c echo.Context) error {
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	batchStandardStudentData := make(map[string]interface{})
	if err := c.Bind(&batchStandardStudentData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	student := &models.Student{ID: uint(newStudentId)}
	if err := student.Find(); err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	batchStandardId := models.GetBatchStandardId(batchStandardStudentData)
	batchStandard := &models.BatchStandard{ID: batchStandardId}
	if err := batchStandard.Find(); err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	err = student.AssignBatchStandard(batchStandard)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrAlreadyHasClass.Error()})
	}
		
	err = student.SaveBalance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "batch Standard Student created", "batch_standard": batchStandard})	
}

func RemoveBatchStandardStudent(c echo.Context) error {

	batchStandardId := c.Param("batch_standard_id")
	newBatchStandardId, err := strconv.Atoi(batchStandardId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	batchStandardStudentId := c.Param("id")
	newBatchStandardStudentId, err := strconv.Atoi(batchStandardStudentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	batchStandardStudent := &models.BatchStandardStudent{ID: uint(newBatchStandardStudentId), BatchStandardId: uint(newBatchStandardId)}
	err = batchStandardStudent.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandardStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}



	err = batchStandardStudent.Delete()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandardStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": newBatchStandardStudentId, "message": "Remove Student successfully"})
}

func UpdateStudentBatchStandard(c echo.Context) error {
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
	if err := bs.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	if err := bs.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Batch Standard updated", "batch_standard": bs})
}

func DeleteStudentBatchStandard(c echo.Context) error {
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

func GetStudentBatchStandards(c echo.Context) error {
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

// func GetBatchUnassignedStandards(c echo.Context) error {
// 	batchID := c.Param("batch_id")
// 	newBatchID, err := strconv.Atoi(batchID)
// 	if err != nil {
// 		fmt.Println("strconv.Atoi failed", err)
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
// 	}	
// 	b := &models.Batch{ID: uint(newBatchID)}
// 	err = b.Find()
// 	if err != nil {
// 		fmt.Println("s.Find(GetBatch)", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
// 	}

// 	bs := &models.BatchStandard{}
// 	standardIds, err := bs.AllIds(b.ID)
// 	if err != nil {
// 		fmt.Println("s.ALL(GetBatchStandardIds)", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
// 	}

// 	standard := &models.Standard{}
// 	standards, err := standard.AllExcept(standardIds)
// 	if err != nil {
// 		fmt.Println("s.ALL(GetBatchStandards)", err)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
// 	}


// 	return c.JSON(http.StatusOK, standards)
// }
