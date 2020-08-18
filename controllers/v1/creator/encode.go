package creator

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/encode"
)

// EncodeProfileList handles listing encode formats
func EncodeProfileList(c echo.Context) error {
	e, err := encode.FormatList()
	if err != nil {
		log.Printf("EncodeProfileList failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, e)
}
