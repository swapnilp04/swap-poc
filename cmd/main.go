package main

import (
	"swapnil-ex/handlers"
	"swapnil-ex/models/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

)

func main() {

	defer db.Close()

	e := echo.New()

	// e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200", "https://labstack.net", "*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.POST("/register", handlers.Register, handlers.OnlySwapnil())
	e.POST("/login", handlers.Login)
	e.PUT("/updateUser", handlers.UpdateUser, handlers.IsLoggedIn)
	e.DELETE("/logout", handlers.Logout, handlers.IsLoggedIn)
	
	e.GET("/students", handlers.GetStudents, handlers.IsLoggedIn)
	e.GET("/students/:id", handlers.GetStudent, handlers.IsLoggedIn)
	e.POST("/students", handlers.CreateStudent, handlers.IsLoggedIn)
	e.PUT("/students/:id", handlers.UpdateStudent, handlers.IsLoggedIn)
	e.DELETE("/students/:id", handlers.DeleteStudent, handlers.IsLoggedIn)
	
	e.GET("/students/:id/hostel", handlers.GetStudentHostel, handlers.IsLoggedIn)
	e.POST("/students/:id/assign_hostel", handlers.AssignStudentHostel, handlers.IsLoggedIn)
	e.PUT("/students/:id/change_hostel", handlers.ChangeStudentHostel, handlers.IsLoggedIn)
	
	e.GET("/students/:student_id/batch_standards", handlers.GetBatchStandardStudents, handlers.IsLoggedIn)
	e.POST("/students/:student_id/batch_standards", handlers.CreateStudentBatchStandard, handlers.IsLoggedIn)

	e.GET("/students/:student_id/transactions", handlers.GetStudentTransactions, handlers.IsLoggedIn)

	e.GET("/standards", handlers.GetStandards, handlers.IsLoggedIn)
	e.GET("/standards/:id", handlers.GetStandard, handlers.IsLoggedIn)
	e.POST("/standards", handlers.CreateStandard, handlers.IsLoggedIn)
	e.PUT("/standards/:id", handlers.UpdateStandard, handlers.IsLoggedIn)
	e.DELETE("/standards/:id", handlers.DeleteStandard, handlers.IsLoggedIn)
	
	e.GET("/batchs", handlers.GetBatchs, handlers.IsLoggedIn)
	e.GET("/batchs/:id", handlers.GetBatch, handlers.IsLoggedIn)
	e.POST("/batchs", handlers.CreateBatch, handlers.IsLoggedIn)
	e.PUT("/batchs/:id", handlers.UpdateBatch, handlers.IsLoggedIn)
	e.DELETE("/batchs/:id", handlers.DeleteBatch, handlers.IsLoggedIn)
	
	e.GET("/batchs/:batch_id/standards", handlers.GetBatchStandards, handlers.IsLoggedIn)
	e.GET("/batchs/:batch_id/unassigned_standards", handlers.GetBatchUnassignedStandards, handlers.IsLoggedIn)
	e.POST("/batchs/:batch_id/batch-standards", handlers.CreateBatchStandard, handlers.IsLoggedIn)
	e.GET("/batchs/:batch_id/batch-standards/:id", handlers.GetBatchStandard, handlers.IsLoggedIn)
	e.PUT("/batchs/:batch_id/batch-standards/:id", handlers.UpdateBatchStandard, handlers.IsLoggedIn)

	e.GET("/hostels", handlers.GetHostels, handlers.IsLoggedIn)
	e.GET("/hostels/:id", handlers.GetHostel, handlers.IsLoggedIn)
	e.POST("/hostels", handlers.CreateHostel, handlers.IsLoggedIn)
	e.PUT("/hostels/:id", handlers.UpdateHostel, handlers.IsLoggedIn)
	e.DELETE("/hostels/:id", handlers.DeleteHostel, handlers.IsLoggedIn)

	e.GET("/hostels/:hostel_id/hostel_rooms", handlers.GetHostelRooms, handlers.IsLoggedIn)
	e.POST("/hostels/:hostel_id/hostel_rooms", handlers.CreateHostelRoom, handlers.IsLoggedIn)
	e.GET("/hostels/:hostel_id/hostel_rooms/:id", handlers.GetHostelRoom, handlers.IsLoggedIn)
	e.PUT("/hostels/:hostel_id/hostel_rooms/:id", handlers.UpdateHostelRoom, handlers.IsLoggedIn)
	e.GET("/hostels/:hostel_id/hostel_rooms/:id/students", handlers.GetHostelRoomStudents, handlers.IsLoggedIn)


	e.Logger.Fatal(e.Start(":8080"))
}
