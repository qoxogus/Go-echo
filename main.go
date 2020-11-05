package main

import (
	"github.com/labstack/echo"
)

func handlelog(c echo.Context) error {
	return c.File("login.html")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", handlelog)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
