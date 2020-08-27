package public

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ListTeams handles listing teams and their members and info
func (r *Repos) ListTeams(c echo.Context) error {
	t, err := r.public.ListTeams(c.Request().Context())
	if err != nil {
		log.Printf("Public ListTeams failed: %+v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, t)
}
