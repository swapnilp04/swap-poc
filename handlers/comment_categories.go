package handlers

import (
	"fmt"
	"net/http"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetCommentCategories(c echo.Context) error {
	// Get all Commnet Categories
	cc := &models.CommentCategory{}
	commentCategories, err := cc.All()
	if err != nil {
		fmt.Println("s.ALL(GetComments)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, commentCategories)
}
