package handlers

import (
	"fmt"
	"net/http"
	//"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetExams(c echo.Context) error {
	// Get all users
	e := &models.Exam{}
	exams, err := e.All()
	if err != nil {
		fmt.Println("s.ALL(GetExams)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, exams)
}
