package main

import (
	"swapnil-ex/handlers"
	"swapnil-ex/models/db"

	"github.com/labstack/echo/v4"
)

func main() {

	defer db.Close()

	e := echo.New()

	// e.Use(middleware.Recover())

	e.POST("/register", handlers.Register, handlers.OnlySwapnil())
	e.POST("/login", handlers.Login)
	e.PUT("/updateUser", handlers.UpdateUser, handlers.IsLoggedIn)
	e.DELETE("/logout", handlers.Logout, handlers.IsLoggedIn)
	e.GET("/students", handlers.GetStudents, handlers.IsLoggedIn)
	e.GET("/students/:id", handlers.GetStudent, handlers.IsLoggedIn)
	e.POST("/students", handlers.CreateStudent, handlers.IsLoggedIn)
	e.PUT("/students/:id", handlers.UpdateStudent, handlers.IsLoggedIn)
	e.DELETE("/students/:id", handlers.DeleteStudent, handlers.IsLoggedIn)

	e.Logger.Fatal(e.Start(":8080"))
}
