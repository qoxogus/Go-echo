package main

import (
	"strings"

	"Scrapper"

	"github.com/labstack/echo"
)

// func main() {
//     scrapper.Scrape("term")
// }

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	term := strings.ToLower(Scrapper.CleanString(c.FormValue("term")))
	Scrapper.Scrape(term)
	// fmt.Println(c.FormValue("term")) //term == python
	// return nil
	return c.Attachment("jobs.csv", "job.csv")
	//return c.File("home.html")
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)  //url에 추가되는 /scrape
	e.Logger.Fatal(e.Start(":1323")) // localhost:1323
}
