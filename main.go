package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

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

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your cat: %#v\n", cat)
	return c.String(http.StatusOK, "we got your cat!")
}

func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed processing addDog request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your dog: %#v", dog)
	return c.String(http.StatusOK, "we got your dog!")
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed processing addHamster request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your hamster: %#v", hamster)
	return c.String(http.StatusOK, "we got your hamster!")
}

// func handleScr(c echo.Context) error {
// 	fmt.Println(c.FormValue("description"))
// 	return nil
// }

func main() {
	fmt.Println("Welcom to the server")
	e := echo.New()
	e.GET("/", handleMain)
	e.GET("/cats/:data", getCats)
	//e.POST("/scrape", _)
	e.POST("/cats/:data", getCats)
	e.Logger.Fatal(e.Start(":1323"))
}

//read code and study
