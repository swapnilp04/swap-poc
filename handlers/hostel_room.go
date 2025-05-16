package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetHostelRooms(c echo.Context) error {
	// Get all users
	
	
	hostel, errHostel := GetHostelForRoom(c)
	if errHostel != nil {
		fmt.Println("s.Find(GetHostelRoom)", errHostel)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	hostelRooms, err := hostel.HostelRooms()
	if err != nil {
		fmt.Println("s.Find(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, hostelRooms)
}

func GetHostelRoom(c echo.Context) error {
	// Get a single user by ID
	hostel, errHostel := GetHostelForRoom(c)

	if errHostel != nil {
		fmt.Println("s.Find(GetHostelRoom)", errHostel)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	hostelRoom := &models.HostelRoom{ID: uint(newId), HostelID: hostel.ID}
	err = hostelRoom.Find()

	if err != nil {
		fmt.Println("s.Find(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, hostelRoom)
}

func GetHostelRoomStudents(c echo.Context) error {
	hostel, errHostel := GetHostelForRoom(c)
	if errHostel != nil {
		fmt.Println("s.Find(GetHostelRoom)", errHostel)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	hostelRoom := &models.HostelRoom{ID: uint(newId), HostelID: hostel.ID}
	err = hostelRoom.Find()

	if err != nil {
		fmt.Println("s.Find(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	hostelRoomStudents, err := hostelRoom.GetHostelRoomStudents()
	if err != nil {
		fmt.Println("s.Find(GetHostelRoomStudents)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, hostelRoomStudents)
}

func CreateHostelRoom(c echo.Context) error {
	hostel, errHostel := GetHostelForRoom(c)
	if errHostel != nil {
		fmt.Println("s.Find(GetHostelRoom)", errHostel)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	hostelRoomData := make(map[string]interface{})

	if err := c.Bind(&hostelRoomData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	fmt.Printf("hostelRooms %+v\n", hostelRoomData)
	hostelRoom := models.NewHostelRoom(hostelRoomData)
	hostelRoom.HostelID = hostel.ID
	hostelRoom.Rate = hostel.Rate
	if err := hostelRoom.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}

	err := hostelRoom.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "hostelRoom created", "hostelRoom": hostelRoom})
}

func UpdateHostelRoom(c echo.Context) error {

	hostel, errHostel := GetHostelForRoom(c)
	if errHostel != nil {
		fmt.Println("s.Find(GetHostelRoom)", errHostel)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	hostelRoomData := make(map[string]interface{})

	if err := c.Bind(&hostelRoomData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	s := &models.HostelRoom{ID: uint(newId)}
	if err := s.Find(); err != nil {
		fmt.Println("s.Find(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	s.Assign(hostelRoomData)
	s.HostelID = hostel.ID
	if err := s.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}
	if err := s.Update(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "HostelRoom updated", "hostelRoom": s})
}

func DeleteHostelRoom(c echo.Context) error {
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	hr := &models.HostelRoom{ID: uint(newId)}
	if err := hr.Find(); err != nil {
		fmt.Println("b.Find(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	if hr.HostelStudentsCount > 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrRemoveHosteStudents.Error()})
	}

	if err := hr.Delete(); err != nil {
		fmt.Println("b.Delete(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	// Delete a user by ID
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "HostelRoom deleted successfully", "id": newId})
}

func RemoveStudentFromHostel(c echo.Context) error {
	id := c.Param("hostel_room_id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	hr := &models.HostelRoom{ID: uint(newId)}
	if err := hr.Find(); err != nil {
		fmt.Println("hr.Find(GetHostelRoom)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	hostelStudentID := c.Param("id")
	newHostelStudentID, err := strconv.Atoi(hostelStudentID)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	err = hr.RemoveHostelRoomStudents(uint(newHostelStudentID))
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Hostel Student deleted successfully", "id": hostelStudentID})
}

func GetHostelForRoom(c echo.Context) (*models.Hostel, error){
	id := c.Param("hostel_id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return &models.Hostel{}, err
	}

	hostel := &models.Hostel{ID: uint(newId)}
	err = hostel.Find()
	return hostel, err
}

