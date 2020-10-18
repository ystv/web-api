package public

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ListTeams handles listing teams and their members and info
//
// @Summary Provides the current teams
// @Description Returns a path of series to a series
// @ID get-public-teams
// @Tags public
// @Produce json
// @Success 200 {array} public.Team
// @Router /v1/public/teams [get]
func (r *Repos) ListTeams(c echo.Context) error {
	t, err := r.public.ListTeams(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("Public ListTeams failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, t)
}
