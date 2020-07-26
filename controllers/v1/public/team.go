package public

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/public"
)

// ListTeams handles listing teams and their members and info
func ListTeams(c echo.Context) error {
	return c.JSON(http.StatusOK, public.ListTeams())
}
