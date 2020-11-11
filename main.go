package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

// func handleMain(c echo.Context) error {
// 	// return c.String(http.StatusOK, "Hello, Taehyeon!!\n")
// 	return c.File("main.html")
// }

func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "yallo from the web side!")
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

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "horay you are one the secret admin main page!")
	//http://localhost:1323/admin/main 으로 들어갔을때 뜨는 문구
}

////////////////////////middleware/////////////////////
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "BlueBot/1.0")
		c.Response().Header().Set("notReallyHeader", "thisHaveNoMeaning")

		return next(c)
	}
}

func main() {
	fmt.Println("Welcom to the server")

	e := echo.New()

	e.Use(ServerHeader)

	g := e.Group("/admin") //middleware를 추가하는 방법 1  그룹에 선언
	// g := e.Group("/admin", middleware.Logger()) //middleware를 추가하는 방법 1  그룹에 선언

	//this logs the server interaction
	// g.Use(middleware.Logger()) //middleware를 추가하는 방법 2  USE로 선언
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// check in the DB    기본인증 middleware
		if subtle.ConstantTimeCompare([]byte(username), []byte("bae")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("0809")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/main", mainAdmin)
	// g.GET("/main", mainAdmin)  //middleware를 추가하는 방법 3  핸들러메서드 뒤 메서드에 직접 추가
	e.GET("/", yallo)
	// e.GET("/", handleMain)
	e.GET("/cats/:data", getCats)

	//e.POST("/scrape", _)
	// e.POST("/cats/:data", getCats)

	e.Logger.Fatal(e.Start(":1323"))
}

//read code and study
