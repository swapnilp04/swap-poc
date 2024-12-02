package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetComments(c echo.Context) error {
	// Get all users
	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}
	
	cm := &models.Comment{}
	comments, err := cm.All(newPage)
	if err != nil {
		fmt.Println("s.ALL(GetComments)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	count, err := cm.AllCount()
	return c.JSON(http.StatusOK, map[string]interface{}{"comments": comments, "total": count})
}

func GetComment(c echo.Context) error {
	// Get a single user by ID
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	comment := &models.Comment{ID: uint(newId)}
	err = comment.Find()
	if err != nil {
		fmt.Println("s.Find(GetBatch)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, comment)
}

func GetStudentComments(c echo.Context) error {
		// Get a single user by ID
	id := c.Param("student_id")
	newId, err := strconv.Atoi(id)
	
	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}

	
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	s := &models.Student{ID: uint(newId)}
	err = s.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	comments, err := s.GetStudentComments(int(newPage))
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	count, err := s.GetStudentCommentsCount()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"comments": comments, "total": count})
}

func GetUpcommingComments(c echo.Context) error {
	// Get all users
	
	cm := &models.Comment{}
	comments, err := cm.UpcommingComments()
	if err != nil {
		fmt.Println("s.ALLUpcommingComments(GetComments)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	return c.JSON(http.StatusOK, comments)
}

func CreateComment(c echo.Context) error {
	cc := c.(CustomContext)
	session := cc.session
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)

	commentData := make(map[string]interface{})
	if err := c.Bind(&commentData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("Comment %+v\n", commentData)
	comment := models.NewComment(commentData)
	comment.UserID = uint(session.UserID)
	comment.StudentID = uint(newStudentId)

	if err := comment.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err = comment.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Comment created", "comment": comment})
}