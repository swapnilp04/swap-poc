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
	e.GET("/students/:id/hostel", handlers.GetStudentHostel, handlers.IsLoggedIn)
	e.POST("/students", handlers.CreateStudent, handlers.IsLoggedIn)
	e.PUT("/students/:id", handlers.UpdateStudent, handlers.IsLoggedIn)
	e.DELETE("/students/:id", handlers.DeleteStudent, handlers.IsLoggedIn)

	e.Logger.Fatal(e.Start(":8080"))
}
