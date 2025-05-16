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

	batchStandard, err := b.GetBatchStandard(uint(newBatchStandardId))
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

	// standardId := c.Param("standard_id")
	// newStandardId, err := strconv.Atoi(standardId)
	// if err != nil {
	// 	fmt.Println("strconv.Atoi failed", err)
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	// }

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

	// standard := &models.Standard{ID: uint(newStandardId)}
	// if err := standard.Find(); err != nil {
	// 	fmt.Println("s.Find(GetBatch)", err)
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	// }

	fmt.Printf("batchs %+v\n", batchStandardData)
	batchStandard := models.NewBatchStandard(batchStandardData, batch)
	// if err := batchStandard.Validate(); err != nil {
	// 	formErr := MarshalFormError(err)	
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	// }

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

func DeactivateBatchStandard(c echo.Context) error {
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

	batchStandard, err := b.GetBatchStandard(uint(newBatchStandardId))
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	err = batchStandard.Deactivate()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchStandard)
}

func ActivateBatchStandard(c echo.Context) error {
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

	batchStandard, err := b.GetBatchStandard(uint(newBatchStandardId))
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	err = batchStandard.Activate()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchStandard)
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

	bs, err := b.GetBatchStandards()
	if err != nil {
		fmt.Println("s.ALL(GetBatchs)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, bs)
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

func GetDefaultBatchStandards(c echo.Context) error {
	
	batch, err := models.GetDefaultBatch()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	batchStandards, err := batch.GetBatchStandards()
	if err != nil {
		fmt.Println("s.ALL(GetBatchStandards)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchStandards)
}

func GetActiveBatchStandards(c echo.Context) error {
	
	batchStandard := &models.BatchStandard{}
	
	batchStandards, err := batchStandard.GetActiveBatchStandards()

	if err != nil {
		fmt.Println("s.ALL(GetActiveBatchStandards)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchStandards)
}

func GetBatchStandardStudents(c echo.Context) error {
	
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
	
	batchStandardStudents, err := 	batchStandard.GetStudents()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	cc := c.(CustomContext)
	if( !(cc.session.Role == "Admin" || cc.session.Role == "Accountant")) {
		masker := mask.NewMasker()
		masker.SetMaskChar("+")
		batchStandardStudents, _ = mask.Mask(batchStandardStudents)
	}

	return c.JSON(http.StatusOK, batchStandardStudents)
}

func GetBatchStandardSubjects(c echo.Context) error {
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
	
	batchStandardSubjects, err := batchStandard.GetSubjects()
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, batchStandardSubjects)
}

func GetBatchStandardSubjectChapters(c echo.Context) error {
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

	subjectId := c.Param("subject_id")
	newSubjectId, err := strconv.Atoi(subjectId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	
	subjectChapters, err := batchStandard.GetChapters(uint(newSubjectId))
	if err != nil {
		fmt.Println("s.Find(GetBatchStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, subjectChapters)
}


func GetBatchStandardLogs(c echo.Context) error {
	
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

	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}
	searchTeacher := c.QueryParam("searchTeacher")
	searchSubject := c.QueryParam("searchSubject")
	searchDate := c.QueryParam("searchDate")
	
	batchStandardLogs, err := batchStandard.GetBatchStandardLogs(newPage, searchTeacher, searchSubject, searchDate)
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	count, err := batchStandard.AllBatchStandardLogsCount(searchTeacher, searchSubject, searchDate)

	return c.JSON(http.StatusOK, map[string]interface{}{"batchStandardLogs": batchStandardLogs, "total": count})
}

func GetBatchStandardReportLogs(c echo.Context) error {
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

	searchDate := c.QueryParam("searchDate")
	batchStandardLogs, err := batchStandard.ReportLogs(searchDate)
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, batchStandardLogs)
}

func GetBatchStandardMonthlyReportLogs(c echo.Context) error {
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

	searchDate := c.QueryParam("searchDate")
	batchStandardLogs, err := batchStandard.ReportMonthlyLogs(searchDate)
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, batchStandardLogs)
}


func GetBatchStandardExams(c echo.Context) error {
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
	
	exams, err := batchStandard.GetExams()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, exams)
}

