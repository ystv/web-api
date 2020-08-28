package stream

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type nginx struct {
	call     string
	addr     string
	clientid string
	app      string
	flashVer string
	swfUrl   string
	tcUrl    string
	pageUrl  string
	name     string
}

// CheckAuth used by nginx to check if the stream has correct credentials
func CheckAuth(c echo.Context) error {
	var allowedKeys [2]string
	allowedKeys[0] = "ystvAreTheBest"
	allowedKeys[1] = "Caspar NRK"
	// log.Printf("Stream name: %s", c.FormValue("name"))
	r := new(nginx)
	err := c.Bind(r)
	// log.Println(c.FormParams())
	if err != nil {
		err = fmt.Errorf("CheckAuth failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if c.FormValue("name") == allowedKeys[0] {
		return c.NoContent(http.StatusCreated)
	}
	return c.NoContent(http.StatusNotFound)
}
