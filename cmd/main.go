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
	e.POST("/students", handlers.CreateStudent, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.PUT("/students/:id", handlers.UpdateStudent, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.DELETE("/students/:id", handlers.DeleteStudent, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	
	e.GET("/students/:id/hostel", handlers.GetStudentHostel, handlers.IsLoggedIn)
	e.POST("/students/:id/assign_hostel", handlers.AssignStudentHostel, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.PUT("/students/:id/change_hostel", handlers.ChangeStudentHostel, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	
	e.GET("/students/:student_id/batch_standards", handlers.GetStudentStandards, handlers.IsLoggedIn)
	e.POST("/students/:student_id/batch_standards", handlers.CreateStudentBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)

	e.GET("/students/:student_id/transactions", handlers.GetStudentTransactions, handlers.IsLoggedIn)
	e.GET("/students/:student_id/transactions/:id", handlers.GetStudentTransaction, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.POST("/students/:student_id/transactions", handlers.PayStudentFee, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.POST("/students/:student_id/transactions/dues/new", handlers.AddStudentDues, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/students/:student_id/transactions/balance", handlers.GetStudentBalance, handlers.IsLoggedIn)

	e.GET("/students/:student_id/student_accounts", handlers.GetStudentAccounts, handlers.IsLoggedIn)
	e.POST("/students/:student_id/student_accounts/deposit", handlers.DepositStudentAccountAmount, handlers.IsLoggedIn)
	e.POST("/students/:student_id/student_accounts/withdraw", handlers.WithdrawStudentAccountAmount, handlers.IsLoggedIn)

	e.GET("/standards", handlers.GetStandards, handlers.IsLoggedIn)
	e.GET("/standards/:id", handlers.GetStandard, handlers.IsLoggedIn)
	e.POST("/standards", handlers.CreateStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/standards/:id", handlers.UpdateStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.DELETE("/standards/:id", handlers.DeleteStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	
	e.GET("/batchs", handlers.GetBatchs, handlers.IsLoggedIn)
	e.GET("/batchs/:id", handlers.GetBatch, handlers.IsLoggedIn)
	e.POST("/batchs", handlers.CreateBatch, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/batchs/:id", handlers.UpdateBatch, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.DELETE("/batchs/:id", handlers.DeleteBatch, handlers.IsLoggedIn, handlers.OnlyAdmin)
	
	e.GET("/batchs/:batch_id/standards", handlers.GetBatchStandards, handlers.IsLoggedIn)
	e.GET("/batchs/:batch_id/unassigned_standards", handlers.GetBatchUnassignedStandards, handlers.IsLoggedIn)
	e.POST("/batchs/:batch_id/batch-standards", handlers.CreateBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/batchs/:batch_id/batch-standards/:id", handlers.GetBatchStandard, handlers.IsLoggedIn)
	e.PUT("/batchs/:batch_id/batch-standards/:id", handlers.UpdateBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/batchs/:batch_id/batch-standards/:id/students", handlers.GetBatchStandardStudents, handlers.IsLoggedIn)

	e.GET("/hostels", handlers.GetHostels, handlers.IsLoggedIn)
	e.GET("/hostels/:id", handlers.GetHostel, handlers.IsLoggedIn)
	e.POST("/hostels", handlers.CreateHostel, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/hostels/:id", handlers.UpdateHostel, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.DELETE("/hostels/:id", handlers.DeleteHostel, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/hostels/:hostel_id/hostel_rooms", handlers.GetHostelRooms, handlers.IsLoggedIn)
	e.POST("/hostels/:hostel_id/hostel_rooms", handlers.CreateHostelRoom, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/hostels/:hostel_id/hostel_rooms/:id", handlers.GetHostelRoom, handlers.IsLoggedIn)
	e.PUT("/hostels/:hostel_id/hostel_rooms/:id", handlers.UpdateHostelRoom, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/hostels/:hostel_id/hostel_rooms/:id/students", handlers.GetHostelRoomStudents, handlers.IsLoggedIn)

	e.GET("/accounts/transactions", handlers.GetTransactions, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/accounts/students/:student_id/transactions/:id", handlers.GetStudentTransaction, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)

	e.GET("/users", handlers.GetUsers, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.POST("/users", handlers.Register, handlers.IsLoggedIn, handlers.OnlyAdmin)


	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
