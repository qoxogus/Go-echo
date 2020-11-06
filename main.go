package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func handleHello(c echo.Context) error {
	return c.String(http.StatusOK, "공부 화이팅")
}

func main() {
	e := echo.New()
	e.GET("/", handleHello)
	e.Logger.Fatal(e.Start(":1323"))
}
