package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func handleMain(c echo.Context) error {
	// return c.String(http.StatusOK, "Hello, Taehyeon!!\n")
	return c.File("main.html")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name") //url중 name=taehyeon 이런식으로 바뀜
	catType := c.QueryParam("type") //url중 name=taehyeon&type=human 이런식으로 바뀜

	dataType := c.Param("data") //여기에 string이나 json이 들어감

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is : %s\nand his type is : %s\n", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to lets us know if you want json or string data",
	})
}

// func handleScr(c echo.Context) error {
// 	fmt.Println(c.FormValue("description"))
// 	return nil
// }

func main() {
	e := echo.New()
	e.GET("/", handleMain)
	e.GET("/cats/:data", getCats)
	//e.POST("/scrape", _)
	e.POST("/cats/:data", getCats)
	e.Logger.Fatal(e.Start(":1323"))
}

//read code and study
