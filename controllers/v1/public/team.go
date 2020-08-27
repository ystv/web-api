package public

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ListTeams handles listing teams and their members and info
func (r *Repos) ListTeams(c echo.Context) error {
	t, err := r.public.ListTeams(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("Public ListTeams failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, t)
}
