package main

import (
	"fmt"

	"github.com/labstack/echo"
)

func handleMain(c echo.Context) error {
	return c.File("main.html")
}

func handleScr(c echo.Context) error {
	fmt.Println(c.FormValue("description"))
	return nil
}

func main() {
	e := echo.New()
	e.GET("/", handleMain)
	e.POST("/scrape", handleScr)
	e.Logger.Fatal(e.Start(":1323"))
}

//read code and study
