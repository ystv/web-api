package clapper

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/clapper/position"
)

// PositionList handles listing all possible positions
func PositionList(c echo.Context) error {
	p, err := position.List()
	if err != nil {
		log.Printf("PositionList failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}
