package main

import (
	"strings"

	"github.com/labstack/echo"
)

// func main() {
// 	scrapper.Scrape("term")
// }

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handlScrape(c echo.Context) error {
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	// fmt.Println(c.FormValue("term")) //term == python
	return nil
	//return c.File("home.html")
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handlScrape)   //url에 추가되는 /scrape
	e.Logger.Fatal(e.Start(":1323")) // localhost:1323
}
