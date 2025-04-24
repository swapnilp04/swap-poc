package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetChapters(c echo.Context) error {

	standard_id := c.Param("standard_id")
	subject_id := c.Param("subject_id")
	subject, err := findStandardSubject(standard_id, subject_id)
	if err != nil {
		fmt.Println("s.Find(GetStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	chapters, err := subject.GetChapters()
	
	if err != nil {
		fmt.Println("s.ALL(GetChapters)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, chapters)
}

func GetChapter(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	chapter := &models.Chapter{ID: uint(newId)}
	err = chapter.Find()
	if err != nil {
		fmt.Println("s.Find(GetChapter)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK,chapter)
}

func CreateChapter(c echo.Context) error {
	standard_id := c.Param("standard_id")
	subject_id := c.Param("subject_id")
	subject, err := findStandardSubject(standard_id, subject_id)
	if err != nil {
		fmt.Println("s.Find(GetStandard)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	newStandard_id, err := strconv.Atoi(standard_id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	chapterData := make(map[string]interface{})
	if err := c.Bind(&chapterData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}
	
	chapter := models.NewChapter(chapterData)
	chapter.SubjectID = subject.ID
	chapter.StandardID = uint(newStandard_id)
	if err := chapter.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err = chapter.Create()
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Chapter created", "chapter": chapter})
}

func UpdateChapter(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	chapterData := make(map[string]interface{})

	if err := c.Bind(&chapterData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	chapter := &models.Chapter{ID: uint(newId)}
	if err := chapter.Find(); err != nil {
		fmt.Println("chapter.Find(GetChapter)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	chapter.Assign(chapterData)
	if err := chapter.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Chapter updated", "chapter": chapter})
}

func DeleteChapter(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	chapter := &models.Chapter{ID: uint(newId)}
	if err := chapter.Find(); err != nil {
		fmt.Println("s.Find(GetChapter)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if err := chapter.Delete(); err != nil {
		fmt.Println("s.Delete(GetChapter)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Chapter deleted successfully"})
}

func findStandardSubject(standard_id string, subject_id string) (models.Subject, error) {
	newStandardId, err := strconv.Atoi(standard_id)
	if err != nil {
		return models.Subject{}, err
	}
	standard := models.Standard{ID: uint(newStandardId)}
	err = standard.Find()

	newSubjectId, err := strconv.Atoi(subject_id)
	if err != nil {
		return models.Subject{}, err
	}
	subject, err := standard.GetSubject(uint(newSubjectId))
	return subject, err
}
