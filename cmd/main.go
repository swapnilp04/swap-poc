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
		AllowOrigins: []string{"http://localhost:4200", "http://eracord.com", "*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.POST("/register", handlers.Register, handlers.OnlySwapnil())
	e.POST("/login", handlers.Login)
	e.PUT("/updateUser", handlers.UpdateUser, handlers.IsLoggedIn)
	e.DELETE("/logout", handlers.Logout, handlers.IsLoggedIn)
	
	e.GET("/students", handlers.GetStudents, handlers.IsLoggedIn)
	e.GET("/students/:id", handlers.GetStudent, handlers.IsLoggedIn)
	e.GET("/students/get-upcomming-birthdays", handlers.GetUpcommingBirthdays, handlers.IsLoggedIn)
	e.POST("/students", handlers.CreateStudent, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.PUT("/students/:id", handlers.UpdateStudent, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.DELETE("/students/:id", handlers.DeleteStudent, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/students/:id/get_exams", handlers.GetStudentExams, handlers.IsLoggedIn)
	e.POST("/students/:id/left_academy", handlers.LeftAcademy, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/students/:id/hostel", handlers.GetStudentHostel, handlers.IsLoggedIn)
	e.POST("/students/:id/assign_hostel", handlers.AssignStudentHostel, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/students/:id/change_hostel", handlers.ChangeStudentHostel, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	
	e.GET("/students/:student_id/batch_standards", handlers.GetStudentStandards, handlers.IsLoggedIn)
	e.POST("/students/:student_id/batch_standards", handlers.CreateStudentBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)

	e.GET("/students/:student_id/transactions", handlers.GetStudentTransactions, handlers.IsLoggedIn)
	e.GET("/students/:student_id/transactions/:id", handlers.GetStudentTransaction, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.POST("/students/:student_id/transactions", handlers.PayStudentFee, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.POST("/students/:student_id/transactions/dues/new", handlers.AddStudentDues, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/students/:student_id/transactions/balance", handlers.GetStudentBalance, handlers.IsLoggedIn)

	e.POST("/students/:student_id/discounts", handlers.AddDiscount, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/students/:student_id/student_accounts", handlers.GetStudentAccounts, handlers.IsLoggedIn)
	e.POST("/students/:student_id/student_accounts/deposit", handlers.DepositStudentAccountAmount, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.POST("/students/:student_id/student_accounts/withdraw", handlers.WithdrawStudentAccountAmount, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/students/:student_id/comments", handlers.GetStudentComments, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)
	e.POST("/students/:student_id/comments", handlers.CreateComment, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)

	e.GET("/standards", handlers.GetStandards, handlers.IsLoggedIn)
	e.GET("/standards/:id", handlers.GetStandard, handlers.IsLoggedIn)
	e.POST("/standards", handlers.CreateStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/standards/:id", handlers.UpdateStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.DELETE("/standards/:id", handlers.DeleteStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/logs", handlers.GetTeacherLogs, handlers.IsLoggedIn)
	e.GET("/logs/:id", handlers.GetTeacherLog, handlers.IsLoggedIn)
	e.POST("/logs", handlers.CreateTeacherLog, handlers.IsLoggedIn, handlers.OnlyTeacher)
	e.PUT("/logs/:id", handlers.UpdateTeacherLog, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)
	e.DELETE("/logs/:id", handlers.DeleteTeacherLog, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)

	e.GET("/log_categories", handlers.GetLogCategories, handlers.IsLoggedIn)
	
	e.GET("/standards/:standard_id/subjects", handlers.GetSubjects, handlers.IsLoggedIn)
	e.GET("/standards/:standard_id/subjects/:id", handlers.GetSubject, handlers.IsLoggedIn)
	e.POST("/standards/:standard_id/subjects", handlers.CreateSubject, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/standards/:standard_id/subjects/:id", handlers.UpdateSubject, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.DELETE("/standards/:standard_id/subjects/:id", handlers.DeleteSubject, handlers.IsLoggedIn, handlers.OnlyAdmin)
	
	e.GET("/teachers", handlers.GetTeachers, handlers.IsLoggedIn)
	e.GET("/teachers/:id", handlers.GetTeacher, handlers.IsLoggedIn)
	e.POST("/teachers", handlers.CreateTeacher, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/teachers/:id", handlers.UpdateTeacher, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/teachers/:id/get_logs", handlers.GetTeachersLog, handlers.IsLoggedIn)

	e.GET("/batchs", handlers.GetBatchs, handlers.IsLoggedIn)
	e.GET("/batchs/:id", handlers.GetBatch, handlers.IsLoggedIn)
	e.POST("/batchs", handlers.CreateBatch, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/batchs/:id", handlers.UpdateBatch, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.DELETE("/batchs/:id", handlers.DeleteBatch, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/batchs/get_default_batch_standards", handlers.GetDefaultBatchStandards, handlers.IsLoggedIn)
	
	e.GET("/batchs/:batch_id/standards", handlers.GetBatchStandards, handlers.IsLoggedIn)
	e.GET("/batchs/:batch_id/unassigned_standards", handlers.GetBatchUnassignedStandards, handlers.IsLoggedIn)
	e.POST("/batchs/:batch_id/batch-standards", handlers.CreateBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/batchs/:batch_id/batch-standards/:id", handlers.GetBatchStandard, handlers.IsLoggedIn)
	e.PUT("/batchs/:batch_id/batch-standards/:id", handlers.UpdateBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/batchs/:batch_id/batch-standards/:id/students", handlers.GetBatchStandardStudents, handlers.IsLoggedIn)
	e.PUT("/batchs/:batch_id/batch-standards/:id/activate", handlers.ActivateBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/batchs/:batch_id/batch-standards/:id/deactivate", handlers.DeactivateBatchStandard, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/batch-standards/:id/subjects", handlers.GetBatchStandardSubjects, handlers.IsLoggedIn)
	e.GET("/batch_standards/:id/get_logs", handlers.GetBatchStandardLogs, handlers.IsLoggedIn)
	e.GET("/batch_standards/:id/get_report_logs", handlers.GetBatchStandardReportLogs, handlers.IsLoggedIn)
	e.GET("/batch_standards/:id/get_monthly_report_logs", handlers.GetBatchStandardMonthlyReportLogs, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.GET("/batch_standards/:id/get_exams", handlers.GetBatchStandardExams, handlers.IsLoggedIn)
	e.GET("/batch-standards/:id/students", handlers.GetBatchStandardStudents, handlers.IsLoggedIn)
	e.DELETE("/batch-standards/:batch_standard_id/batch-standard-students/:id", handlers.RemoveBatchStandardStudent, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/hostels", handlers.GetHostels, handlers.IsLoggedIn)
	e.GET("/hostels/:id", handlers.GetHostel, handlers.IsLoggedIn)
	e.GET("/hostels/get_early_expired_students", handlers.GetEarlyExpiredHostelStudents, handlers.IsLoggedIn)
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
	e.GET("/accounts/transactions/report", handlers.GetTransactionsReport, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/users", handlers.GetUsers, handlers.IsLoggedIn)
	e.GET("/users/current", handlers.GetCurrentUser, handlers.IsLoggedIn)
	e.POST("/users", handlers.Register, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.PUT("/users/update_password", handlers.UpdatePassword, handlers.IsLoggedIn)
	e.POST("/users/:id/deactive_user", handlers.DeactivateUser, handlers.IsLoggedIn, handlers.OnlyAdmin)
	e.POST("/users/:id/active_user", handlers.ActivateUser, handlers.IsLoggedIn, handlers.OnlyAdmin)

	e.GET("/comments", handlers.GetComments, handlers.IsLoggedIn)
	e.GET("/comments/:id", handlers.GetComment, handlers.IsLoggedIn)
	e.POST("/comments/:id/completed_comment", handlers.CommentCompleted, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.GET("/comments/upcomming_comments", handlers.GetUpcommingComments, handlers.IsLoggedIn)

	e.GET("/comment_categories", handlers.GetCommentCategories, handlers.IsLoggedIn)

	e.GET("/exams", handlers.GetExams, handlers.IsLoggedIn)
	e.GET("/exams/:id", handlers.GetExam, handlers.IsLoggedIn)
	e.POST("/exams", handlers.CreateExam, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)
	e.PUT("/exams/:id", handlers.UpdateExam, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.DELETE("/exams/:id", handlers.DeleteExam, handlers.IsLoggedIn, handlers.OnlyAdminAccountant)
	e.POST("/exams/:id/conduct_exam", handlers.ConductExam, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)
	e.POST("/exams/:id/publish_exam", handlers.PublishExam, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)
	e.GET("/exams/:id/exam_students", handlers.GetExamStudents, handlers.IsLoggedIn)
	e.POST("/exams/:id/save_exam_marks", handlers.SaveExamMarks, handlers.IsLoggedIn, handlers.OnlyAdminAccountantClerkTeacher)
	e.GET("/exams/get_exam_group_report", handlers.GetExamGroupReport, handlers.IsLoggedIn)


	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
